package controller

//controller/login.go 登入用

import (
	"backendAPI/response"

	. "backendAPI/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Chatboteventid_de struct {
	EventID   string `json:"eventID"`
	AutoID    string `json:"autoID"`
	Keyword   string `json:"keyword"`   //關鍵字 list,json
	Msg       string `json:"msg"`       //自動回覆內容
	Status    string `json:"status"`    //0:不啟用 1:啟用
	Weekday   string `json:"weekday"`   //1byte,8個bit 星期一到日
	StatTime  string `json:"statTime"`  //每日起始時間(分鐘為單位)0-1440
	EndTime   string `json:"endTime"`   //每日結束時間(分鐘為單位)(須注意跨日問題)
	StartDate string `json:"startDate"` //起始日期
	EndDate   string `json:"endDate"`   //結束日期

}

func Delchatbot(c *gin.Context) *response.Response {

	autoID := StrToInt(c.PostForm("autoID")) //活動ID

	//insert into autoresponse(eventID,keyword,msg,status,weekday,startTime,endTime,startDate,endDate) values(?,?,?,?,?,?,?,?,?)
	d_qs := "DELETE FROM autoresponse WHERE autoID = ?"
	if del := DB.Exec(d_qs, autoID); del.Error != nil {
		Log.Error("[Delchatbot] ", del.Error)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
