package controller

import (
	. "chatbotAPI/lib"
	"chatbotAPI/model"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
)

func ReceiveMQ() {
	mq := GetMQ()
	fmt.Println("ReceiveMQ->")
	ch, err := mq.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"janus-events", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message:%s\n", d.Body)

			var mapData []map[string]interface{}
			json.Unmarshal(d.Body, &mapData)
			code := int(mapData[0]["type"].(float64))

			var textroom = 64
			if code != textroom {
				Log.Errorf("type is not textroom, code:%s", code)
				return
			}

			var event model.Event
			err := mapstructure.Decode(mapData[0]["event"], &event)
			if err != nil {
				Log.Errorf("decode err:%s", err.Error())
				return
			}

			what := event.Data["event"].(string)
			fmt.Println("what: ", what)

			switch what {
			case "message":
				InsertMessageData([]byte(event.Data["text"].(string)))

				TestAIMessageDate([]byte(event.Data["text"].(string)))
			case "join", "leave":
				UpdateLastTime(int64(mapData[0]["timestamp"].(float64)), int(event.Data["room"].(float64)), event.Data["username"].(string), event.Data["display"].(string))

				// UpdateLastTime(int64(mapData[0]["timestamp"].(float64)), int(event.Data["room"].(float64)), event.Data["username"].(string), "")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//寫入訊息資料
func InsertMessageData(message []byte) {
	influxDB := GetInflux()
	defer influxDB.Close()

	var msg model.Message
	// json.NewDecoder(strings.NewReader(c.PostForm("message"))).Decode(&msg)

	var mapData map[string]interface{}
	json.Unmarshal(message, &mapData)

	janusMsg := mapData["janusMsg"]

	if janusMsg != nil {
		err := mapstructure.Decode(janusMsg, &msg)
		if err != nil {
			Log.Errorf("[Insert] decode err:%s", err.Error())
			return
		}
		// fmt.Println("msg.Message: ", msg.MsgContent)
	} else {
		Log.Warn("[Insert] 資料格式錯誤!")
		return
	}

	if len(msg.ChatroomID) != 7 {
		Log.Warn("[Insert] 資料格式錯誤")
		return
	}

	Log.Infof("[Insert] janusMsg: %s", janusMsg)

	tags := map[string]string{
		"chatroomID": msg.ChatroomID,
	}

	for_jsonStr, _ := json.Marshal(msg.Format)

	con_jsonStr, _ := json.Marshal(msg.Config)

	fields := map[string]interface{}{
		"sender":  int(msg.Sender),
		"type":    int(msg.Type),
		"msgType": int(*msg.MsgType),
		"message": msg.MsgContent,
		"format":  string(for_jsonStr),
		"config":  string(con_jsonStr),
		// "format":  *msg.Format,
	}

	insErr := InsertInflux(influxDB, "message", tags, fields)
	if insErr != nil {
		Log.Errorf("[Insert] err:%s", insErr.Error)
		return
	}
	Log.Info("Insert Ok")
}

//更新在線時間
func UpdateLastTime(unix int64, eventid int, chatroomid string, admin string) {
	var u_qs string
	switch admin {
	case "user":
		u_qs = "UPDATE chatroom SET lastVisitB = ? WHERE eventID = ? AND chatroomID = ?"
		time := time.Unix(0, int64(unix)*int64(time.Microsecond))
		ins := DB.Exec(u_qs, time, eventid, chatroomid)
		if ins.Error != nil {
			Log.Errorf("[Update] err:%s", ins.Error)
			return
		}
	case "admin":
		u_qs = "UPDATE chatroom SET lastVisit = ? WHERE eventID = ?"
		time := time.Unix(0, int64(unix)*int64(time.Microsecond))
		ins := DB.Exec(u_qs, time, eventid)
		if ins.Error != nil {
			Log.Errorf("[Update] err:%s", ins.Error)
			return
		}
	}
	Log.Info("Update Ok")
	return
}
