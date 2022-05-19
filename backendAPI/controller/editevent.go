package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//編輯活動
func EditEvent(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	eventId := context.Param("eventId")
	Log.Infof("[EditEvent] eventId:%s", eventId)

	tx := DB.Begin()

	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM event WHERE eventID = ? LIMIT 1)"
	if err := tx.Raw(s_qs, eventId).Row().Scan(&exists); err != nil {
		tx.Rollback()
		Log.Error("[EditEvent] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有此活動"}, http.StatusNotFound)
	}

	u_qs := "UPDATE event SET"
	keys := []string{"name", "description", "callable", "homeurl", "message", "status", "icon"}
	var args []interface{}
	var qParts []string
	for _, key := range keys {
		value := context.PostForm(key)
		if value != "" {
			qParts = append(qParts, fmt.Sprintf(" %s = ?", key))
			args = append(args, value)
		}
	}

	if args != nil {
		u_qs += strings.Join(qParts, ",") + " WHERE eventID = ?"
		args = append(args, eventId)
		// fmt.Printf("u_qs: %s\n", u_qs)

		ins := tx.Exec(u_qs, args...)
		if ins.Error != nil {
			tx.Rollback()
			Log.Error("[EditEvent] ", ins.Error)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
	}

	//更新客服人員
	cslist := context.PostForm("cslist")
	if !IsEmpty(cslist) {
		var founder int
		e_qs := "SELECT accountID FROM event WHERE eventID=?"
		if err := tx.Raw(e_qs, eventId).Row().Scan(&founder); err != nil {
			tx.Rollback()
			Log.Error("[EditEvent] ", err)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}

		d_qs := "DELETE FROM accountEvent WHERE accountID!= ? AND eventID = ?"
		if del := tx.Exec(d_qs, founder, eventId); del.Error != nil {
			tx.Rollback()
			Log.Error("[EditEvent] ", del.Error)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}

		c := []byte(cslist)

		var list []map[string]interface{}
		json.Unmarshal([]byte(c), &list)

		// var acts []map[string]interface{}

		for _, r := range list {
			act_id := int(r["accountID"].(float64))
			if act_id == founder {
				continue
			}
			if act_id == 0 {
				// if len(acts) <= 0 {
				// 	resp, err := AccountInfo(context.GetString("token"))
				// 	if err != nil {
				// 		Log.Errorf("[AccountInfo] %s", err)
				// 	}
				// 	if len(resp) <= 0 {
				// 		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
				// 	}
				// 	acts = resp
				// }

				// roleInfo := GetRoleInfo(acts, int(r["userNo"].(float64)))

				// act := roleInfo["LoginName"].(string)
				// n := roleInfo["UserName"].(string)

				//新增帳號
				name := r["name"].(string)
				nickname := r["nickname"].(string)
				if IsEmpty(nickname) || IsEmpty(name) {
					return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
				}

				n_id, err := AddAccount(context.GetInt("customer"), Viewer, name, nickname)
				if err != nil {
					tx.Rollback()
					Log.Errorf("[EditAdmin] %s", err)
					return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
				}
				act_id = n_id
			}

			ae_qs := "INSERT INTO accountEvent(accountID, eventID) VALUES(?,?)"
			if ins := tx.Exec(ae_qs, act_id, eventId); ins.Error != nil {
				tx.Rollback()
				Log.Error("[EditEvent] ", ins.Error)
				return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
			}
		}
	}

	tx.Commit()

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
