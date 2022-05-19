package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//取得使用者資訊
func GetUserInfo(context *gin.Context) *response.Response {
	chatroomId := context.Param("chatroomId")
	Log.Infof("[GetUserInfo] chatroomId:%s", chatroomId)

	//檢查此帳號是否有聊天室權限
	exists := IsChatRoom(context.GetInt("admin"), context.GetInt("id"), chatroomId)
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	var (
		mobile      int
		name        string
		icon        int
		lastTime    string
		tag         json.RawMessage
		description string
	)

	qs := `SELECT IFNULL(mobile, ""), IFNULL(userNameB, ""), icon, lastVisitB, IFNULL(tag,"[]"), IFNULL(descriptionB,"") FROM chatroom WHERE chatroomID = ?`
	// qs := `SELECT IFNULL(mobile, ""), IFNULL(userNameB, ""), icon, tag, IFNULL(descriptionB,"") FROM chatroom WHERE chatroomID = ?`
	err := DB.Raw(qs, chatroomId).Row().Scan(
		&mobile,
		&name,
		&icon,
		&lastTime,
		&tag,
		&description,
	)

	if err != nil {
		Log.Error("[GetUserInfo] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	t, err := time.Parse(time.RFC3339, lastTime)
	if err != nil {
		Log.Error("[GetUserInfo] ", err)
	}
	unixTime := t.UnixNano()

	return response.Resp().Json(gin.H{"mobile": mobile, "name": name, "icon": icon, "tag": tag, "description": description, "lastVisit": unixTime}, http.StatusOK)
}
