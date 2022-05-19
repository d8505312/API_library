package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type often struct {
	SmsID   int    `json:"smsID"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

//取得常用簡訊
func GetOftenSMS(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	accountId := context.Param("accountId")
	Log.Infof("[GetOftenSMS] accountId:%s", accountId)

	var o often
	var oftens []often

	qs := `SELECT smsID, subject, content FROM oftenSMS WHERE accountID = ?`
	rows, err := DB.Raw(qs, accountId).Rows()
	if err != nil {
		Log.Error("[GetOftenSMS] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&o.SmsID,
			&o.Subject,
			&o.Content,
		)
		oftens = append(oftens, o)
	}

	if oftens == nil {
		oftens = make([]often, 0)
	}

	return response.Resp().Json(gin.H{"oftens": oftens}, http.StatusOK)
}
