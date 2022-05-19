package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//新增活動
func AddEvent(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	accountID := context.PostForm("accountID")
	name := context.PostForm("name")
	description := context.PostForm("description")
	callable := context.PostForm("callable")
	homeurl := context.PostForm("homeurl")
	messagelist := context.PostForm("messageList")
	icon := context.PostForm("icon")
	cslist := context.PostForm("cslist")

	Log.Infof("[AddEvent] accountID:%s | name:%s | description:%s | callable:%s | homeurl:%s | messagelist:%s | icon:%s | cslist:%s", accountID, name, description, callable, homeurl, messagelist, icon, cslist)

	tx := DB.Begin()
	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM account WHERE accountID = ?)"
	if err := tx.Raw(s_qs, &accountID).Row().Scan(&exists); err != nil || !exists {
		tx.Rollback()
		Log.Error("[AddEvent] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	//insert event table
	i_qs := "INSERT INTO event(accountID, name, description, callable, homeurl, message, icon) VALUES(?,?,?,?,?,?,?)"
	ins := tx.Exec(i_qs, &accountID, &name, &description, &callable, &homeurl, &messagelist, &icon)
	if ins.Error != nil {
		tx.Rollback()
		Log.Error("[AddEvent] ", ins.Error)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	// get last insert id
	var id uint64
	err := tx.Raw("SELECT LAST_INSERT_ID()").Row().Scan(&id)
	if err != nil {
		tx.Rollback()
		Log.Error("[AddEvent] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	// fmt.Println("LAST NISERT　eventID: ", id)
	Log.Info("[AddEvent] LAST_INSERT_ID:", id)

	//insert accountEvent table
	c := []byte(cslist)

	var list []map[string]interface{}
	json.Unmarshal([]byte(c), &list)
	myid, _ := strconv.ParseFloat(accountID, 64)
	list = append(list, map[string]interface{}{"accountID": myid})

	for _, r := range list {
		act_id := int(r["accountID"].(float64))
		if act_id == 0 {
			//新增帳號
			name := r["name"].(string)
			nickname := r["nickname"].(string)
			if IsEmpty(name) || IsEmpty(nickname) {
				return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
			}
			n_id, err := AddAccount(context.GetInt("customer"), Viewer, name, nickname)
			if err != nil {
				tx.Rollback()
				Log.Errorf("[AddEvent] %s", err)
				return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
			}
			act_id = n_id
		}
		ae_qs := "INSERT INTO accountEvent(accountID, eventID) VALUES(?,?)"
		if ins := tx.Exec(ae_qs, act_id, &id); ins.Error != nil {
			tx.Rollback()
			Log.Error("[AddEvent] ", ins.Error)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
	}

	tx.Commit()
	return response.Resp().Json(gin.H{"id": id}, http.StatusOK)
}
