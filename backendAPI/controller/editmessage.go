package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//編輯訊息
func EditMessage(context *gin.Context) *response.Response {
	chatroomId := context.PostForm("chatroomID")
	msgType := context.PostForm("msgType")
	id := context.PostForm("id")
	config := context.PostForm("config")

	Log.Infof("[EditMessage] chatroomId:%s | msgType:%s | id:%s | config:%s\n", chatroomId, msgType, id, config)

	if IsEmpty(chatroomId) || IsEmpty(msgType) || IsEmpty(id) || IsEmpty(config) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	influx := GetInflux()
	defer influx.Close()

	sender := 0 // 0:後台 1:前台

	qs := fmt.Sprintf("SELECT message, format, type FROM message WHERE sender = %d AND chatroomID = '%s' AND config =~ /%s/", sender, chatroomId, id)

	resp, err := QueryInflux(influx, qs)
	if err != nil {
		Log.Error("[EditMessage] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	// fmt.Println("resp", resp)
	if len(resp[0].Series) == 0 {
		return response.Resp().Json(gin.H{"msg": "訊息不存在"}, http.StatusNotFound)
	}

	time, _ := time.Parse(time.RFC3339, fmt.Sprint(resp[0].Series[0].Values[0][0]))
	message := fmt.Sprint(resp[0].Series[0].Values[0][1])
	format := fmt.Sprint(resp[0].Series[0].Values[0][2])
	types := fmt.Sprint(resp[0].Series[0].Values[0][3])

	tags := map[string]string{
		"chatroomID": chatroomId,
	}

	fields := map[string]interface{}{
		"sender":  sender,
		"type":    StrToInt(types),
		"msgType": StrToInt(msgType),
		"message": message,
		"format":  format,
		"config":  config,
	}

	insErr := InsertInflux(influx, "message", tags, fields, time)
	if insErr != nil {
		Log.Errorf("[Insert] err:%s", insErr.Error)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
