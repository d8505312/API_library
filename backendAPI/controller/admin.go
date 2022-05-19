package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type admin struct {
	AccountID int32  `json:"accountID"`
	Name      string `json:"name"`
}

//取得管理人員列表
func GetAdminList(context *gin.Context) *response.Response {
	if context.GetInt("admin") != Manager {
		// return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	var a admin
	var admins []admin
	qs := "SELECT accountID, account FROM account WHERE customerID = ? AND admin = ?"
	rows, err := DB.Raw(qs, context.GetInt("customer"), Editor).Rows()
	if err != nil {
		Log.Error("[CustomerList] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&a.AccountID,
			&a.Name,
		)
		admins = append(admins, a)
	}

	if admins == nil {
		return response.Resp().Json(gin.H{"msg": "沒有管理人員"}, http.StatusNotFound)
	}

	return response.Resp().Json(gin.H{"admins": admins}, http.StatusOK)
}
