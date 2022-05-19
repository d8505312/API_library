package controller

//controller/login.go 登入用

import (
	"backendAPI/response"
	"strconv"
	"strings"
	"time"

	. "backendAPI/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type Message struct {
// 	ChatroomID string                 `json:"chatroomID" binding:"required"`
// 	Sender     int                    `json:"sender" binding:"required"`
// 	Type       int                    `json:"type" binding:"required"`
// 	MsgType    *int                   `json:"msgType" binding:"required"`
// 	MsgContent string                 `json:"msgContent" binding:"required"`
// 	Time       int64                  `json:"Time" binding:"required"`
// 	Format     map[string]interface{} `json:"format" binding:"required"`
// 	Config     map[string]interface{} `json:"config" binding:"required"`
// }

func Timeconvert(unix_s string, timetype string) string {
	var time_s string
	unix, err := strconv.ParseInt(unix_s, 10, 64)
	if err != nil {
		Log.Error("[error] 轉換失敗", err)

	}
	janus_date := time.Unix(0, unix) //janus來源的時間
	switch timetype {
	case "h":

		hmst := janus_date.Hour()*60 + janus_date.Minute()
		// Log.Info("[hmst格式]", hmst)
		time_s = strconv.Itoa(hmst)
	case "d":
		time_s = janus_date.Format("2006-01-02")

	}
	// Log.Info("[time格式]", unix_s, unix, timetype)
	// Log.Info("[輸出的time格式]", time_s)
	return time_s
}

func Addchatbot(c *gin.Context) *response.Response {
	// var msg_js Message
	// var mapData map[string]interface{}
	eventID := c.PostForm("eventID") //活動ID
	subject := c.PostForm("subject") //標題
	keyword := c.PostForm("keyword") //關鍵字listjson格式
	msg := c.PostForm("msg")         //自動回覆內容
	status := c.PostForm("status")   //0:不啟用 ,1:啟用
	weekday := c.PostForm("weekday") // 1byte,8個bit分別代表代表星期一到日,高bitalways0
	weekday_arry := strings.Split(weekday, "")

	startTime := Timeconvert(c.PostForm("startTime"), "h") //每日起試時間(分鐘為單位)0-1440
	endTime := Timeconvert(c.PostForm("endTime"), "h")     //每日結束時間(分鐘為單位)(須注意跨日問題)0-1440
	startDate := Timeconvert(c.PostForm("startDate"), "d") //起始日期
	endDate := Timeconvert(c.PostForm("endDate"), "d")     //結束日期
	i := 0
	total := 0
	// var msg_data []byte = []byte(msg) //將進來的質轉成byte
	// json.Unmarshal(msg_data, &mapData)
	// janusMsg := mapData["janusMsg"]
	// // Log.Info("解析情況：", janusMsg)
	// if janusMsg != nil {
	// 	err := mapstructure.Decode(janusMsg, &msg_js)
	// 	if err != nil {
	// 		Log.Error("[Addchatbot] decode err:%s", err.Error())

	// 	}

	// }
	for {
		// Log.Info("讓我看看", weekday_arry[i], weekday_arry, i)
		switch weekday_arry[i] {
		case "1":
			total = total + 1
			// Log.Info("1.加數1：", total)
		case "2":
			total = total + 2
			// Log.Info("2.加數2：", total)
		case "3":
			total = total + 4
			// Log.Info("3.加數4：", total)
		case "4":
			total = total + 8
			// Log.Info("4.加數8：", total)
		case "5":
			total = total + 16
			// Log.Info("5.加數16：", total)
		case "6":
			total = total + 32
			// Log.Info("6.加數32：", total)
		case "7":
			total = total + 64
			// Log.Info("7.加數64：", total)
		}
		i = i + 1
		if i == len(weekday_arry) {
			break
		}

	}
	// msg_ins := ""
	// switch *msg_js.MsgType {
	// case 1:
	// 	msg_ins = msg_js.MsgContent
	// case 6:
	// 	fmt.Println("msg目前沒圖片")
	// case 7:
	// 	fmt.Println("msg目前沒檔案")
	// }
	// Log.Info("加總數", total)
	Log.Info("[Addchatbot] eventID:%s | subject:%s | keyword:%s | msg:%s | status:%s | weekday:%s | startTime:%s | endTime:%s | startDate:%s | endDate:%s", eventID, &subject, &keyword, msg, status, total, startTime, endTime, startDate, endDate)

	tx := DB.Exec("insert into autoresponse(eventID,subject,keyword,msg,status,weekday,startTime,endTime,startDate,endDate) values(?,?,?,?,?,?,?,?,?,?)", &eventID, &subject, &keyword, msg, &status, &total, &startTime, &endTime, &startDate, &endDate)

	if tx.Error != nil {
		Log.Error("[Addchatbot] insert err:%s", tx.Error)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	return response.Resp().Json(gin.H{"eventID": eventID, "subject": subject, "keyword": keyword, "msg": msg, "status": status, "weekdaytotal": total, "startTime": startTime, "endTime": endTime, "startDate": startDate, "endDate": endDate}, http.StatusOK)
}
