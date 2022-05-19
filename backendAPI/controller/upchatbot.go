package controller

//controller/login.go 登入用

import (
	"backendAPI/response"
	"strings"

	. "backendAPI/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Upchatbot(c *gin.Context) *response.Response {

	eventID := c.PostForm("eventID") //活動ID
	autoID := c.PostForm("autoID")   //活動ID
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
	// Log.Info("加總數", total)
	Log.Info("[Upchatbot] eventID:%v | autoID:%v | subject:%s | keyword:%v | msg:%v | status:%v | weekday:%v | startTime:%v | endTime:%v | startDate:%v | endDate:%v", eventID, autoID, keyword, msg, status, total, startTime, endTime, startDate, endDate)
	//update user set name=? where userID=?
	//insert into autoresponse(eventID,keyword,msg,status,weekday,startTime,endTime,startDate,endDate) values(?,?,?,?,?,?,?,?,?)
	tx := DB.Exec("update autoresponse set subject=?, keyword=?,msg=?,status=?,weekday=?,startTime=?,endTime=?,startDate=?,endDate=? where eventID=? and autoID= ?", &subject, &keyword, &msg, &status, &total, &startTime, &endTime, &startDate, &endDate, &eventID, &autoID)

	if tx.Error != nil {

		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
