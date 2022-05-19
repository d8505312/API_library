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

type chatroom struct {
	ChatroomID string `json:"chatroomID"`
	LastVisit  string `json:"-"`
	Icon       int    `json:"icon"`
	Name       string `json:"name"`
	Mobile     int    `json:"mobile"`
	Time       int64  `json:"time"`
	Msg        string `json:"msg"`
	MsgType    int    `json:"msgType"`
	PinTop     int8   `json:"pinTop"` //0:否 1:是
	Unread     int8   `json:"unread"` //0:無 1:有
}

//聊天室列表
func ChatRoomList(context *gin.Context) *response.Response {
	eventId := context.Param("eventId")
	Log.Infof("[ChatRoomList] eventId:%s", eventId)

	//檢查是否有活動
	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM event WHERE eventID = ? LIMIT 1)"
	if err := DB.Raw(s_qs, eventId).Row().Scan(&exists); err != nil {
		Log.Error("[ChatRoomList] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	if !exists {
		return response.Resp().Json(gin.H{"msg": "不存在任何聊天室"}, http.StatusNotFound)
	}

	//檢查是否有參與此活動
	exists = IsEvent(context.GetInt("admin"), context.GetInt("id"), StrToInt(eventId))
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	influx := GetInflux()
	defer influx.Close()

	//查詢活動內的所有聊天室id
	var c chatroom
	var chatrooms []chatroom
	qs := "SELECT chatroomID, lastVisit, mobile, icon, userNameB, pinTop FROM chatroom WHERE eventID = ?"
	rows, err := DB.Raw(qs, eventId).Rows()

	if err != nil {
		Log.Error("[ChatRoomList] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&c.ChatroomID,
			&c.LastVisit,
			&c.Mobile,
			&c.Icon,
			&c.Name,
			&c.PinTop,
		)

		//查詢聊天室的最新一則訊息
		m_qs := fmt.Sprintf("SELECT message, msgType FROM message WHERE chatroomID = '%s' AND msgType < %d AND config =~ /%s/ ORDER BY DESC LIMIT 1", c.ChatroomID, 10, `"recallStatus":false`)

		resp, err := QueryInflux(influx, m_qs)
		if err != nil {
			Log.Error("[ChatRoomList] ", err)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
		c.Unread = 0
		if len(resp[0].Series) == 0 {
			c.Msg = ""
			c.MsgType = 1
		} else {
			values := resp[0].Series[0].Values[0]
			c.Msg = fmt.Sprint(values[1])
			c.MsgType, _ = strconv.Atoi(fmt.Sprint(values[2]))

			//檢查訊息是否已讀
			t1, _ := time.Parse(time.RFC3339, c.LastVisit)
			t2, _ := time.Parse(time.RFC3339, fmt.Sprint(values[0]))
			c.Time = t2.UnixNano()

			// fmt.Printf("t1-> %s\n t2-> %s\n", t1, t2)
			if t1.Before(t2) {
				c.Unread = 1
			}
		}
		chatrooms = append(chatrooms, c)
	}

	if chatrooms == nil {
		return response.Resp().Json(gin.H{"msg": "不存在任何聊天室"}, http.StatusNotFound)
	}

	return response.Resp().Json(gin.H{"chatroomList": chatrooms}, http.StatusOK)
}
