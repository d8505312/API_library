package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

//新增常用簡訊
func AddOftenSMS(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	accountID := context.PostForm("accountID")
	subject := context.DefaultPostForm("subject", "")
	content := context.PostForm("content")

	Log.Infof("[AddOftenSMS] accountID:%s | subject:%s | content:%s", accountID, subject, content)

	if IsEmpty(accountID) || IsEmpty(content) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	tx := DB.Begin()
	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM account WHERE accountID = ?)"
	if err := tx.Raw(s_qs, &accountID).Row().Scan(&exists); err != nil || !exists {
		tx.Rollback()
		Log.Error("[AddOftenSMS] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	//insert oftenSMS table
	i_qs := "INSERT INTO oftenSMS(accountID, subject, content) VALUES(?,?,?)"
	ins := tx.Exec(i_qs, &accountID, &subject, &content)
	if ins.Error != nil {
		tx.Rollback()
		Log.Error("[AddOftenSMS] ", ins.Error)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	tx.Commit()
	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
