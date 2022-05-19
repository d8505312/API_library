package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type eventCustomer struct {
	AccountID int32  `json:"accountID"`
	Name      string `json:"name"`
}

//取得活動客服人員列表
func GetEventCustomerList(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	eventId := context.Param("eventId")
	Log.Infof("[GetEventCustomerList] eventId:%s", eventId)

	var cs eventCustomer
	var csList []eventCustomer
	qs := "SELECT a.accountID, a.account FROM  account a, accountEvent ae WHERE a.accountID = ae.accountID AND ae.eventID = ? AND a.admin < ?"
	rows, err := DB.Raw(qs, eventId, Manager).Rows()
	if err != nil {
		Log.Error("[GetEventCustomerList] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&cs.AccountID,
			&cs.Name,
		)
		csList = append(csList, cs)
	}

	if csList == nil {
		return response.Resp().Json(gin.H{"msg": "沒有客服人員"}, http.StatusNotFound)
	}

	return response.Resp().Json(gin.H{"cslist": csList}, http.StatusOK)
}

type Customer struct {
	AccountID int32 `json:"accountID"`
}

//更新活動客服人員列表
func PutEventCustomerList(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	eventId := context.Param("eventId")
	Log.Infof("[PutEventCustomerList] eventId:%s", eventId)
	// Debug info
	// fmt.Println("c.Request.Method >> " + context.Request.Method)
	// fmt.Println("c.Request.URL.String() >> " + context.Request.URL.String())

	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM event WHERE eventID = ? LIMIT 1)"
	if err := DB.Raw(s_qs, eventId).Row().Scan(&exists); err != nil {
		Log.Error("[PutEventCustomerList] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有此活動"}, http.StatusNotFound)
	}

	d_qs := "DELETE FROM accountEvent WHERE eventID = ?"
	if del := DB.Exec(d_qs, eventId); del.Error != nil {
		Log.Error("[PutEventCustomerList] ", del.Error)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	var cslist []Customer
	if err := context.ShouldBindJSON(&cslist); err != nil {
		Log.Error("[PutEventCustomerList] ", err)
		return response.Resp().Json(gin.H{"msg": err.Error()}, http.StatusBadRequest)
	}
	for _, r := range cslist {
		fmt.Println("AccountID: ", r.AccountID)
		i_qs := "INSERT INTO accountEvent(accountID, eventID) VALUES(?,?)"
		if ins := DB.Exec(i_qs, r.AccountID, eventId); ins.Error != nil {
			Log.Error("[PutEventCustomerList] ", ins.Error)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
	}

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
