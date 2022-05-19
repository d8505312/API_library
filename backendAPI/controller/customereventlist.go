package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type customerEvent struct {
	EventID int32  `json:"eventID"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
}

//取得客服人員活動列表
func GetCustomerEventList(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	accountId := context.Param("accountId")
	Log.Infof("[GetCustomerEventList] accountId:%s", accountId)

	var accountName string
	s_qs := "SELECT account FROM account WHERE accountID = ?"
	if err := DB.Raw(s_qs, accountId).Row().Scan(&accountName); err != nil {
		Log.Error("[GetCustomerEventList] ", err)
		return response.Resp().Json(gin.H{"msg": "帳號不存在"}, http.StatusNotFound)
	}

	var e customerEvent
	var events []customerEvent
	qs := "SELECT e.eventID, e.name, e.icon FROM  event e, accountEvent ae WHERE e.eventID = ae.eventID AND ae.accountID = ? AND e.status < ?"
	rows, err := DB.Raw(qs, accountId, Manager).Rows()
	if err != nil {
		Log.Error("[GetCustomerEventList] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&e.EventID,
			&e.Name,
			&e.Icon,
		)
		events = append(events, e)
	}

	// if events == nil {
	// 	return response.Resp().Json(gin.H{"msg": "沒有活動"}, http.StatusNotFound)
	// }

	return response.Resp().Json(gin.H{"accountName": accountName, "events": events}, http.StatusOK)
}

type Event struct {
	EventID int32 `json:"eventID"`
}

//更新客服人員活動列表
func PutCustomerEventList(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	accountId := context.Param("accountId")
	Log.Infof("[PutCustomerEventList] accountId:%s", accountId)

	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM account WHERE accountID = ? LIMIT 1)"
	if err := DB.Raw(s_qs, accountId).Row().Scan(&exists); err != nil {
		Log.Error("[PutCustomerEventList] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	if !exists {
		return response.Resp().Json(gin.H{"msg": "帳號不存在"}, http.StatusNotFound)
	}

	d_qs := "DELETE FROM accountEvent WHERE accountID = ?"
	if del := DB.Exec(d_qs, accountId); del.Error != nil {
		Log.Error("[PutCustomerEventList] ", del.Error)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	var events []Event
	if err := context.ShouldBindJSON(&events); err != nil {
		Log.Error("[PutCustomerEventList] ", err)
		return response.Resp().Json(gin.H{"msg": err.Error()}, http.StatusBadRequest)
	}
	for _, r := range events {
		fmt.Println("EventID: ", r.EventID)
		i_qs := "INSERT INTO accountEvent(accountID, eventID) VALUES(?,?)"
		if ins := DB.Exec(i_qs, accountId, r.EventID); ins.Error != nil {
			Log.Error("[PutCustomerEventList] ", ins.Error)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
	}

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
