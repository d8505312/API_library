package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type livechat struct {
	ChatroomID string    `json:"chatroomID"`
	LastVisitB string    `json:"-"`
	Msg        string    `json:"msg"`
	MsgType    int       `json:"msgType"`
	Unread     int8      `json:"unread"` //0:無 1:有
	Time       time.Time `json:"-"`
}

//活動中最近有訊息的聊天室
func LiveChatRoom(context *gin.Context) *response.Response {
	eventId := context.Param("eventId")
	Log.Infof("[LiveChatRoom] eventId:%s", eventId)

	//檢查是否有活動
	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM event WHERE eventID = ? LIMIT 1)"
	if err := DB.Raw(s_qs, eventId).Row().Scan(&exists); err != nil {
		Log.Error("[LiveChatRoom] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有任何訊息"}, http.StatusNotFound)
	}

	//檢查是否有此活動權限
	exists = IsEvent(context.GetInt("admin"), context.GetInt("id"), StrToInt(eventId))
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	influx := GetInflux()
	defer influx.Close()

	//查詢活動內的所有聊天室id
	var l livechat
	var livechats []livechat
	qs := "SELECT chatroomID, lastVisit FROM chatroom WHERE eventID = ?"
	rows, err := DB.Raw(qs, eventId).Rows()
	if err != nil {
		Log.Error("[LiveChatRoom] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&l.ChatroomID,
			&l.LastVisitB,
		)

		//查詢聊天室的最新一則訊息
		m_qs := fmt.Sprintf("SELECT message, msgType FROM message WHERE chatroomID = '%s' AND msgType < %d  ORDER BY DESC LIMIT 1", l.ChatroomID, 10)

		resp, err := QueryInflux(influx, m_qs)
		if err != nil {
			Log.Error("[LiveChatRoom] ", err)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
		l.Unread = 0
		if len(resp[0].Series) == 0 {
			l.Msg = "沒有聊天訊息"
			l.MsgType = 1
		} else {
			values := resp[0].Series[0].Values[0]
			l.Msg = fmt.Sprint(values[1])
			l.MsgType, _ = strconv.Atoi(fmt.Sprint(values[2]))

			//檢查訊息是否已讀
			lastVisit, _ := time.Parse(time.RFC3339, l.LastVisitB)
			l.Time, _ = time.Parse(time.RFC3339, fmt.Sprint(values[0]))

			// fmt.Printf("lastVisit-> %s\n l.Time-> %s\n", lastVisit, l.Time)
			if lastVisit.Before(l.Time) {
				l.Unread = 1
			}
		}
		livechats = append(livechats, l)
	}

	if livechats == nil {
		return response.Resp().Json(gin.H{"msg": "沒有任何訊息"}, http.StatusNotFound)
	}

	//依時間排序篩選最近十筆資料
	sort.Slice(livechats, func(i, j int) bool {
		return livechats[i].Time.After(livechats[j].Time)
	})

	livechats = livechats[:Min(len(livechats), 10)]

	return response.Resp().Json(gin.H{"chatroomList": livechats}, http.StatusOK)
}
