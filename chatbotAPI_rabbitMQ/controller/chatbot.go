package controller

import (
	. "chatbotAPI/lib"
	"chatbotAPI/model"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

type SendreplytrMQ struct {
	ChatroomID string      `json:"chatroomID"`
	MsgType    int         `json:"msgType"`
	Sender     int         `json:"sender"`
	MsgContent string      `json:"msgContent"`
	Time       int64       `json:"time"`
	Types      int         `json:"type"`
	Format     Formattype  `json:"format"`
	Configs    ConfigsType `json:"config"`
}
type ConfigsType struct {
	Ids               string `json:"id"`
	IsReply           bool   `json:"isReply"`
	ReplyObj          string `json:"replyObj"`
	CurrentDate       string `json:"currentDate"`
	IsExpire          bool   `json:"isExpire"`
	IsPlay            bool   `json:"isPlay"`
	IsRead            bool   `json:"isRead"`
	MsgFunctionStatus bool   `json:"msgFunctionStatus"`
	MsgMoreStatus     bool   `json:"msgMoreStatus"`
	RecallPopUp       bool   `json:"recallPopUp"`
	RecallStatus      bool   `json:"recallStatus"`
}
type Formattype struct {
}

//****時間屬性轉換
func timedate(time_layout, time_date string) time.Time {

	timehmsloc, err := time.ParseInLocation(time_layout, time_date, time.Local) //調整時區
	if err != nil {
		fmt.Println(err)
	}
	timehmsunix := timehmsloc.Unix()                            //纳秒時間戳
	timehmspar := time.Unix(timehmsunix, 0).Format(time_layout) //取得想要的部分
	timehmspar_t, err := time.Parse(time_layout, timehmspar)    //轉換成時間
	if err != nil {
		fmt.Println(err)
	}
	return timehmspar_t
}

//****時間判斷
func timejudge(traget, time_s, time_e, saetype int) bool {
	var judge bool
	switch saetype {
	case 0: //無換日
		if time_s <= traget && traget <= time_e {

			judge = true

		} else {
			judge = false
		}
	case 1: //開始日 (開始大於結束的換日)
		if time_s <= traget && traget <= 1440 {
			judge = true
		} else {
			judge = false
		}
	case 2: //結束日(開始大於結束的換日)
		if 0 <= traget && traget <= time_e {
			judge = true
		} else {
			judge = false
		}
	case 3: //非開始結束日(開始大於結束的換日)
		if time_s <= traget || traget <= time_e {
			judge = true
		} else {
			judge = false
		}
	}
	return judge
}

//****判斷是否有符合list清單資料
func listjudge(listdata []string, data_traget string) bool {
	list_int := 0
	for i := 0; i < len(listdata); i++ {
		if strings.Contains(data_traget, listdata[i]) {
			list_int = list_int + 1

		}
	}
	if list_int > 0 {
		return true
	} else {
		return false
	}

}

//****不同時間組合判斷 ex:week+date+Hms or week + date ....
func timecombo(week, ymd, hms bool) bool {
	if week && ymd && hms {
		return true
	} else {
		return false
	}

}
func TestAIMessageDate(message []byte) {

	//寫入訊息資料
	var msg model.Message
	var mapData map[string]interface{}
	var types string
	var ymd_bool, hms_bool, week_bool bool
	ymd_bool, hms_bool, week_bool = true, true, true
	json.Unmarshal(message, &mapData)

	janusMsg := mapData["janusMsg"]

	if janusMsg != nil {
		err := mapstructure.Decode(janusMsg, &msg)
		if err != nil {
			Log.Errorf("[Insert] decode err:%s", err.Error())
			return
		}
	} else {
		Log.Warn("[Insert] 資料格式錯誤!")
		return
	}
	if len(msg.ChatroomID) != 7 {
		Log.Warn("[Insert] 資料格式錯誤")
	}
	if *msg.MsgType != 1 {
		Log.Warn("[Insert] 非文字訊息")
	} else {
		// inputReader := bufio.NewReader(os.Stdin) // 從輸入讀取資料

		message_input := msg.MsgContent
		//****拆成 星期、日期、時間
		janus_date := time.Unix(0, msg.Time) //janus來源的時間
		yt, mt, dt := janus_date.Date()
		ymd := yt*10000 + int(mt)*100 + dt
		week_jans := janus_date.Weekday()
		hmst := janus_date.Hour()*60 + janus_date.Minute()

		Log.Infof("[Show] roomid: %v", msg.ChatroomID)
		types = "keywords"                          //關鍵字、時間、回覆內容 等等類別名稱
		keystr, _ := Getdata(msg.ChatroomID, types) //取得鍵質內容
		//****開始日期
		timeymd_s, _ := Getdata(msg.ChatroomID, "dates") //取值
		timeymd_st := timedate("2006/1/2", timeymd_s)    //時間轉換
		y_st, m_st, d_st := timeymd_st.Date()
		ymd_st := y_st*10000 + int(m_st)*100 + d_st

		//****結束日期
		timeymd_e, _ := Getdata(msg.ChatroomID, "datee") //取得鍵質內容
		timeymd_et := timedate("2006/1/2", timeymd_e)    //時間轉換
		y_et, m_et, d_et := timeymd_et.Date()
		ymd_et := y_et*10000 + int(m_et)*100 + d_et

		//****開始時間
		timehms_s, _ := Getdata(msg.ChatroomID, "times") //取得鍵質內容
		hms_st, _ := strconv.Atoi(timehms_s)

		//****結束時間
		timehms_e, _ := Getdata(msg.ChatroomID, "timee") //取得鍵質內容
		hms_et, _ := strconv.Atoi(timehms_e)

		//****星期
		week, _ := Getdata(msg.ChatroomID, "week") //取得鍵質內容

		//****關鍵字
		types = "reply"                              //關鍵字、時間、回覆內容 等等類別名稱
		replytr, _ := Getdata(msg.ChatroomID, types) //取得鍵質內容

		// fmt.Printf("歡迎(%v) 回來", msg.ChatroomID)

		keystrlist := strings.Split(keystr, ",") //關鍵字應用
		weeklist := strings.Split(week, "")
		week_str := strconv.Itoa(int(week_jans))
		if week != "" {
			week_bool = listjudge(weeklist, week_str)
			Log.Infof("星期判斷 %v", week_bool)
		}
		if timeymd_s != "" && timeymd_e != "" {
			ymd_bool = timejudge(ymd, ymd_st, ymd_et, 0)
			Log.Infof("日期判斷 %v", ymd_bool)
		}
		if timehms_s != "" && timehms_e != "" {
			if hms_st < hms_et {
				hms_bool = timejudge(hmst, hms_st, hms_et, 0)
				Log.Infof("時間判斷 %v", hms_bool)
			}
			//新增換日條間 排除第一天跟最後一天
			if ymd == ymd_st && hms_st > hms_et { //若為開始日
				hms_bool = timejudge(hmst, hms_st, hms_et, 1)
			}
			if ymd == ymd_et && hms_st > hms_et { //若為結束日
				hms_bool = timejudge(hmst, hms_st, hms_et, 2)
			}
			if ymd != ymd_st && ymd != ymd_et && hms_st > hms_et { //非開始跟結束日
				hms_bool = timejudge(hmst, hms_st, hms_et, 3)
			}
			// fmt.Println("條件時間 :", ymd, ymd_st, ymd_et, hmst, hms_st, hms_et)
			// os.Exit(1)
		}

		if listjudge(keystrlist, message_input) {

			if timecombo(week_bool, ymd_bool, hms_bool) {

				send := SendreplytrMQ{
					ChatroomID: msg.ChatroomID,
					MsgType:    1,
					Sender:     0, // 0:客服,1:使用者
					MsgContent: replytr,
					Time:       time.Now().Unix(),
					Types:      1,
					Format:     Formattype{},
					Configs: ConfigsType{
						Ids:               "1230",
						IsReply:           false,
						ReplyObj:          "",
						CurrentDate:       time.Now().Format("2006-01-02"),
						IsExpire:          false,
						IsPlay:            false,
						IsRead:            false,
						MsgFunctionStatus: false,
						MsgMoreStatus:     false,
						RecallPopUp:       false,
						RecallStatus:      false,
					},
				}

				sendBytes, err := json.Marshal(send)
				if err != nil {
					Log.Errorf("[Insert] err:%s", err)
				}
				SendToMQ(sendBytes)
				Log.Infof("送入時間 %v 送入時間 %v", janus_date, hms_st)
				fmt.Println(replytr)

			} else {
				Log.Infof("不在規定時間 %v", janus_date)

			}

		} else {
			Log.Infof("非關鍵字 %v", message_input)
		}
	}
}
