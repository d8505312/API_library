package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type searchMsg struct {
	Time   int64  `json:"time"`   //訊息時間
	Sender int    `json:"sender"` //發送者
	Text   string `json:"text"`   //訊息內容
}

//搜尋訊息
func SearchMsg(context *gin.Context) *response.Response {
	chatroomId := context.PostForm("chatroomID")
	keyword := context.PostForm("keyword")
	Log.Infof("[SearchMsg] chatroomId:%s | keyword:%s", chatroomId, keyword)

	if IsEmpty(chatroomId) || IsEmpty(keyword) {
		return response.Resp().Json(gin.H{"msg": "聊天室不存在"}, http.StatusNotFound)
	}

	//檢查此帳號是否有聊天室權限
	exists := IsChatRoom(context.GetInt("admin"), context.GetInt("id"), chatroomId)
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	influx := GetInflux()
	defer influx.Close()

	qs := fmt.Sprintf("SELECT sender,message FROM message WHERE  chatroomID = '%s' AND msgType < %d AND message =~ /%s/", chatroomId, 10, keyword)

	resp, err := QueryInflux(influx, qs)
	if err != nil {
		Log.Error("[SearchMsg] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	//fmt.Println("resp", resp)
	if len(resp[0].Series) == 0 {
		return response.Resp().Json(gin.H{"msg": "找不到符合的訊息"}, http.StatusNotFound)
	}

	var s searchMsg
	var list []searchMsg

	for _, row := range resp[0].Series[0].Values {
		t, _ := time.Parse(time.RFC3339, fmt.Sprint(row[0]))
		s.Time = t.Unix()
		s.Sender, _ = strconv.Atoi(fmt.Sprint(row[1]))
		s.Text = fmt.Sprint(row[2])
		list = append(list, s)
	}

	return response.Resp().Json(gin.H{"searchList": list}, http.StatusOK)
}
