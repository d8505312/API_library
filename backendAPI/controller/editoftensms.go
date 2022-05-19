package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//編輯常用簡訊
func EditOftenSMS(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	smsID := context.PostForm("smsID")
	Log.Infof("[EditOftenSMS] smsID:%s", smsID)

	if IsEmpty(smsID) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	var exists bool
	s_qs := "SELECT EXISTS(SELECT smsID FROM oftenSMS WHERE smsID = ? LIMIT 1)"
	if err := DB.Raw(s_qs, smsID).Row().Scan(&exists); err != nil {
		Log.Error("[EditOftenSMS] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	if !exists {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	u_qs := "UPDATE oftenSMS SET"
	keys := []string{"subject", "content"}
	var args []interface{}
	var qParts []string
	for _, key := range keys {
		value := context.PostForm(key)
		if !IsEmpty(value) {
			qParts = append(qParts, fmt.Sprintf(" %s = ?", key))
			args = append(args, value)
		}
	}
	if args != nil {
		u_qs += strings.Join(qParts, ",") + " WHERE smsID = ?"
		args = append(args, smsID)
		// fmt.Printf("u_qs: %s\n", u_qs)

		ins := DB.Exec(u_qs, args...)
		if ins.Error != nil {
			Log.Error("[EditOftenSMS] ", ins.Error)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
	}

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
