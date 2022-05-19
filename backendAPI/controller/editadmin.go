package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

//編輯管理人員列表
func EditAdmin(context *gin.Context) *response.Response {
	if context.GetInt("admin") != Manager {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	accounts := context.PostForm("accounts")
	admin := context.PostForm("admin")
	Log.Infof("[EditAdminList] accounts:%s | admin:%s", accounts, admin)

	if IsEmpty(admin) || IsEmpty(accounts) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}
	tx := DB.Begin()

	a := []byte(accounts)

	var list []map[string]interface{}
	json.Unmarshal([]byte(a), &list)

	for _, r := range list {
		var exists bool
		act_id := int(r["accountID"].(float64))
		if act_id == 0 {
			//新增帳號
			name := r["name"].(string)
			nickname := r["nickname"].(string)
			if IsEmpty(nickname) || IsEmpty(name) {
				return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
			}
			_, err := AddAccount(context.GetInt("customer"), StrToInt(admin), name, nickname)
			if err != nil {
				Log.Errorf("[EditAdmin] %s", err)
				return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
			}
		} else {
			s_qs := "SELECT EXISTS(SELECT accountID FROM account WHERE customerID = (SELECT customerID FROM account WHERE accountID =?) AND accountID = ? LIMIT 1)"
			if err := tx.Raw(s_qs, context.GetInt("id"), act_id).Row().Scan(&exists); err != nil {
				tx.Rollback()
				Log.Error("[EditAdminList] ", err)
				return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
			}

			if !exists {
				tx.Rollback()
				return response.Resp().Json(gin.H{"msg": "帳號不存在"}, http.StatusNotFound)
			}

			u_qs := "UPDATE account SET admin = ? WHERE accountID = ?"
			ins := tx.Exec(u_qs, admin, act_id)
			if ins.Error != nil {
				tx.Rollback()
				Log.Error("[EditAdminList] ", ins.Error)
				return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
			}
		}
	}

	tx.Commit()

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
