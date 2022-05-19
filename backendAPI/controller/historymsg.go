package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type janusMsg struct {
	JanusMsg historyMsg `json:"janusMsg"`
}

type historyMsg struct {
	MsgContent string          `json:"msgContent"`
	Format     json.RawMessage `json:"format"`
	Config     json.RawMessage `json:"config"`
	Type       int             `json:"type"`
	MsgType    int             `json:"msgType"`
	Sender     int             `json:"sender"`
	Time       int64           `json:"time"`
}

//聊天室過往訊息
func HistoryMsg(context *gin.Context) *response.Response {
	chatroomId := context.PostForm("chatroomID")

	Log.Infof("[HistoryMsg] chatroomId:%s", chatroomId)

	if IsEmpty(chatroomId) {
		return response.Resp().Json(gin.H{"msg": "聊天室不存在"}, http.StatusNotFound)
	}
	start := context.DefaultPostForm("start", "0")
	size := context.DefaultPostForm("size", "20")

	//檢查此帳號是否有聊天室權限
	exists := IsChatRoom(context.GetInt("admin"), context.GetInt("id"), chatroomId)
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	influx := GetInflux()
	defer influx.Close()

	qs := fmt.Sprintf("SELECT message, format, msgType, sender, config, type FROM message WHERE  chatroomID = '%s' AND msgType < %d ORDER BY DESC LIMIT %s OFFSET %s", chatroomId, 10, size, start)

	resp, err := QueryInflux(influx, qs)
	if err != nil {
		Log.Error("[HistoryMsg] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	//fmt.Println("resp", resp)

	var h janusMsg
	var list []janusMsg

	if len(resp[0].Series) == 0 {
		return response.Resp().Json(gin.H{"totalCount": 0, "messageList": list}, http.StatusOK)
		// wecomeMsg, err := getWecomeMsg(chatroomId)
		// if err != nil {
		// 	Log.Error("[getWecomeMsg] ", err)
		// 	return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		// }
		// json.Unmarshal(wecomeMsg, &list)
		// return response.Resp().Json(gin.H{"totalCount": len(list), "messageList": list}, http.StatusOK)
	}

	values := resp[0].Series[0].Values
	for _, row := range values {
		t, _ := time.Parse(time.RFC3339, fmt.Sprint(row[0]))
		h.JanusMsg.Time = t.UnixNano()
		h.JanusMsg.MsgContent = fmt.Sprint(row[1])
		formatStr := fmt.Sprint(row[2])
		if strings.Compare(formatStr, "<nil>") == 0 || strings.Compare(formatStr, "json") == 0 || strings.Compare(formatStr, "") == 0 {
			h.JanusMsg.Format = nil
		} else {
			h.JanusMsg.Format = json.RawMessage(formatStr)
		}
		h.JanusMsg.MsgType, _ = strconv.Atoi(fmt.Sprint(row[3]))
		h.JanusMsg.Sender, _ = strconv.Atoi(fmt.Sprint(row[4]))

		configStr := fmt.Sprint(row[5])
		if strings.Compare(configStr, "<nil>") == 0 || strings.Compare(formatStr, "") == 0 {
			h.JanusMsg.Config = nil
		} else {
			h.JanusMsg.Config = json.RawMessage(configStr)
		}
		h.JanusMsg.Type, _ = strconv.Atoi(fmt.Sprint(row[6]))
		list = append(list, h)
	}

	return response.Resp().Json(gin.H{"totalCount": len(values), "messageList": list}, http.StatusOK)
}

// func getWecomeMsg(chatroomId string) ([]byte, error) {
// 	var wecomeMsg json.RawMessage
// 	qs := "SELECT message FROM event WHERE eventID = (SELECT eventID FROM chatroom WHERE chatroomID = ?)"
// 	if err := DB.Raw(qs, chatroomId).Row().Scan(&wecomeMsg); err != nil {
// 		return nil, err
// 	}
// 	return wecomeMsg, nil
// }
