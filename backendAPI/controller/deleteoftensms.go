package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

//刪除常用簡訊
func DeleteOftenSMS(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	smsID := context.PostForm("smsID")

	Log.Infof("[DeleteOftenSMS] smsID:%s", smsID)

	if IsEmpty(smsID) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	tx := DB.Begin()
	var exists bool
	s_qs := "SELECT EXISTS(SELECT smsID FROM oftenSMS WHERE smsID = ?)"
	if err := tx.Raw(s_qs, &smsID).Row().Scan(&exists); err != nil || !exists {
		tx.Rollback()
		Log.Error("[DeleteOftenSMS] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	d_qs := "DELETE FROM oftenSMS WHERE smsID = ?"
	if del := DB.Exec(d_qs, &smsID); del.Error != nil {
		Log.Error("[DeleteOftenSMS] ", del.Error)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	tx.Commit()
	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
