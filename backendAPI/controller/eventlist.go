package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type event struct {
	EventID int32  `json:"eventID"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Status  int32  `json:"status"` //0:停用 1:啟用 2:刪除
	Unread  int8   `json:"unread"` //0:無 1:有
}

//活動資訊列表
func Eventlist(context *gin.Context) *response.Response {
	admin := context.GetInt("admin")
	var e event
	var events []event
	var qs string
	var args []interface{}

	switch admin {
	case Viewer, Editor:
		qs = "SELECT e.eventID, e.name, e.icon, e.status FROM event e, accountEvent ae WHERE e.eventID = ae.eventID AND ae.accountID = ? AND e.status < ?"
		args = append(args, context.GetInt("id"), 2)
	case Manager:
		qs = "SELECT e.eventID, e.name, e.icon, e.status FROM event e,account a WHERE e.accountID = a.accountID AND a.customerID=? AND e.status < ?"
		args = append(args, context.GetInt("customer"), 2)
	}

	rows, err := DB.Raw(qs, args...).Rows()
	if err != nil {
		Log.Error("[Eventlist] rows:", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&e.EventID,
			&e.Name,
			&e.Icon,
			&e.Status,
		)

		e.Unread, err = isUnread(int(e.EventID))
		if err != nil {
			Log.Errorf("[isUnread] %s\n", err)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}

		events = append(events, e)
	}

	if events == nil {
		return response.Resp().Json(gin.H{"msg": "沒有活動"}, http.StatusNotFound)
	}

	return response.Resp().Json(gin.H{"eventList": events}, http.StatusOK)
}

func isUnread(eventId int) (unread int8, err error) {
	influx := GetInflux()
	defer influx.Close()

	var chatroomID, lastTime string
	qs := "SELECT chatroomID, lastVisit FROM chatroom WHERE eventID = ?"
	rows, err := DB.Raw(qs, eventId).Rows()
	if err != nil {
		return unread, fmt.Errorf("error: rows is %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&chatroomID,
			&lastTime,
		)

		t, _ := time.Parse(time.RFC3339, lastTime)
		last := t.UnixNano()

		m_qs := fmt.Sprintf("SELECT sender FROM message WHERE sender = %d AND chatroomID = '%s' AND time > %d-2s ORDER BY DESC LIMIT 1", 1, chatroomID, last)

		resp, err := QueryInflux(influx, m_qs)
		if err != nil {
			return unread, fmt.Errorf("error: resp is %s", err)
		}

		if len(resp[0].Series) > 0 {
			unread = 1
			return unread, nil
		}
	}
	return unread, nil
}
