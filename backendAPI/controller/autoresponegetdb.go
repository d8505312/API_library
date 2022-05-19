package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type autoresponeeventid struct {
	EventID   string `json:"eventID"`
	AutoID    string `json:"autoID"`
	Keyword   string `json:"keyword"`   //關鍵字 list,json
	Msg       string `json:"msg"`       //自動回覆內容
	Status    string `json:"status"`    //0:不啟用 1:啟用
	Weekday   string `json:"weekday"`   //1byte,8個bit 星期一到日
	StartTime string `json:"statTime"`  //每日起始時間(分鐘為單位)0-1440
	EndTime   string `json:"endTime"`   //每日結束時間(分鐘為單位)(須注意跨日問題)
	StartDate string `json:"startDate"` //起始日期
	EndDate   string `json:"endDate"`   //結束日期

}

//活動規則列表
func Autoresponse(context *gin.Context) *response.Response {
	var bot Chatboteventid
	var chatbots []Chatboteventid
	var qs string
	eventID := context.Param("eventId")
	autoID := context.Param("autoId")
	fmt.Print("[Getdbchatbot] eventID:%s| autoID:%s", eventID, autoID)
	if autoID == "0" {

		qs = "SELECT eventID,autoID,keyword,msg,status,weekday,startTime,endTime,startDate,endDate FROM autoresponse  where eventID = ? or autoID = ?" //

	} else {

		qs = "SELECT eventID,autoID,keyword,msg,status,weekday,startTime,endTime,startDate,endDate FROM autoresponse  where eventID = ? and autoID = ?"

	}
	rows, err := DB.Raw(qs, eventID, autoID).Rows()
	if err != nil {
		Log.Error("[Getdbchatbot] rows:", err)
		return response.Resp().Json(gin.H{"msg": err}, http.StatusInternalServerError)
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&bot.EventID,
			&bot.AutoID,
			&bot.Keyword,
			&bot.Msg,
			&bot.Status,
			&bot.Weekday,
			&bot.StartTime,
			&bot.EndTime,
			&bot.StartDate,
			&bot.EndDate,
		)
		Weekdaytotla := StrToInt(bot.Weekday)

		var Weekday_en string
		if Weekdaytotla == 0 {
			Weekday_en = ""
		} else if Weekdaytotla == 127 {
			Weekday_en = "1234567"
		} else {
			result := strings.Split(fmt.Sprintf("%07b", Weekdaytotla), "")

			for k, v := range result {
				if v == "1" {
					Weekday_en = Weekday_en + fmt.Sprint(k+1)
				}
			}
		}
		bot.Weekday = Weekday_en
		bot.StartTime = Timechange(bot.StartTime, "m")
		bot.EndTime = Timechange(bot.EndTime, "m")
		bot.StartDate = Timechange(bot.StartDate, "d")
		bot.EndDate = Timechange(bot.EndDate, "d")
		fmt.Printf("[開始時間]%s|[結束時間]%s|[開始日期]%s|[結束日期]%s|", bot.StartTime, bot.EndTime, bot.StartDate, bot.EndDate)
		// Log.Info("時間輸出：", bot.StartTime, bot.EndTime, bot.StartDate, bot.EndDate)
		Log.Info("[Getdbchatbot]：", bot.EventID, bot.AutoID, bot.Keyword, bot.Status, bot.Weekday, bot.StartTime, bot.EndTime, bot.StartDate, bot.EndDate)
		chatbots = append(chatbots, bot)
	}
	if chatbots == nil {
		return response.Resp().Json(gin.H{"msg": "找不到chatbot"}, http.StatusNotFound)
	}

	return response.Resp().Json(gin.H{"chatbotrule": chatbots}, http.StatusOK)
}
