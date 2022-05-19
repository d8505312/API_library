package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//編輯使用者資訊
func EditUserInfo(context *gin.Context) *response.Response {
	chatroomId := context.PostForm("chatroomID")
	Log.Infof("[EditUserInfo] chatroomId:%s", chatroomId)

	if IsEmpty(chatroomId) {
		return response.Resp().Json(gin.H{"msg": "聊天室不存在"}, http.StatusNotFound)
	}

	var exists bool
	s_qs := "SELECT EXISTS(SELECT chatroomID FROM chatroom WHERE chatroomID = ? LIMIT 1)"
	if err := DB.Raw(s_qs, chatroomId).Row().Scan(&exists); err != nil {
		Log.Error("[EditUserInfo] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	if !exists {
		return response.Resp().Json(gin.H{"msg": "聊天室不存在"}, http.StatusNotFound)
	}

	u_qs := "UPDATE chatroom SET"
	keys := []string{"name", "mobile", "icon", "tag", "description", "pinTop"}
	var args []interface{}
	var qParts []string
	for _, key := range keys {
		value := context.PostForm(key)
		if value != "" {
			if key == "name" {
				key = "userNameB"
			}
			if key == "description" {
				key = "descriptionB"
			}
			qParts = append(qParts, fmt.Sprintf(" %s = ?", key))
			args = append(args, value)
		}
	}

	if args != nil {
		u_qs += strings.Join(qParts, ",") + " WHERE chatroomID = ?"
		args = append(args, chatroomId)
		// fmt.Printf("u_qs: %s\n", u_qs)

		ins := DB.Exec(u_qs, args...)
		if ins.Error != nil {
			Log.Error("[EditUserInfo] ", ins.Error)
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
	}

	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
