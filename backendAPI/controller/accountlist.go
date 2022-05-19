package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accounts struct {
	AccountID int    `json:"accountID"`
	Name      string `json:"name"`
	Admin     int    `json:"admin"`
	Nickname  string `json:"nickname,omitempty"`
	// UserNo int    `json:"userNo"`
}

//所有子帳號清單
func AccountList(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	token := context.GetString("token")

	if !CheckConnection(token) {
		return response.Resp().Json(gin.H{"msg": "token無效,請重新登入"}, http.StatusUnauthorized)
	}

	respData, err := AccountInfo(token)
	if err != nil {
		Log.Errorf("[AccountInfo] %s", err)
	}

	if len(respData) <= 0 {
		return response.Resp().Json(gin.H{"msg": "沒有帳號存在"}, http.StatusNotFound)
	}

	tx := DB.Begin()

	var a accounts
	var accounts []accounts

	for _, d := range respData {
		if int(d["RoleID"].(float64)) != 1000 {
			account := d["LoginName"].(string)
			name := d["UserName"].(string)
			// a.UserNo = int(d["UserNo"].(float64))

			var exists bool
			s_qs := "SELECT EXISTS(SELECT accountID FROM account WHERE customerID = (SELECT customerID FROM account WHERE accountID =?) AND account = ? AND admin < ? LIMIT 1)"
			if err := tx.Raw(s_qs, context.GetInt("id"), account, 2).Row().Scan(&exists); err != nil {
				tx.Rollback()
				Log.Error("[AccountList] ", err)
				return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
			}

			var accountId, admin int
			if exists {
				s_qs := "SELECT accountID,admin FROM account WHERE account = ? AND admin < ?"
				if err := tx.Raw(s_qs, account, 2).Row().Scan(&accountId, &admin); err != nil {
					tx.Rollback()
					Log.Error("[AccountList] ", err)
					return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
				}
				a.AccountID = accountId
				a.Name = account
				a.Admin = admin
			} else {
				a.AccountID = 0
				a.Admin = 0
				a.Name = account
				a.Nickname = name
			}
			accounts = append(accounts, a)
		}
	}

	tx.Commit()

	return response.Resp().Json(gin.H{"accounts": accounts}, http.StatusOK)
}
