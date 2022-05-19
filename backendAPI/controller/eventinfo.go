package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"net/http"

	"gopkg.in/guregu/null.v4"

	"github.com/gin-gonic/gin"
)

type info struct {
	Name        null.String     `json:"name"`
	Description null.String     `json:"description"`
	Icon        null.String     `json:"icon"`
	Callable    int32           `json:"callable"`    //0:false 1:true
	Jid         int32           `json:"jid"`         //Janus ServerID
	Homeurl     null.String     `json:"homeurl"`     //活動官網 url
	Message     json.RawMessage `json:"messageList"` //歡迎訊息
	Status      int32           `json:"status"`      //0:停用 1:啟用 2:刪除
}

//活動資訊
func EventInfo(context *gin.Context) *response.Response {
	eventId := context.Param("eventId")
	Log.Infof("[EventInfo] eventId:%s", eventId)

	//檢查是否有參與此活動
	exists := IsEvent(context.GetInt("admin"), context.GetInt("id"), StrToInt(eventId))
	if !exists {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	var info info
	qs := "SELECT name, description, icon, callable, homeurl, message, status FROM event WHERE eventID =?"
	err := DB.Raw(qs, eventId).Row().Scan(
		&info.Name,
		&info.Description,
		&info.Icon,
		&info.Callable,
		&info.Homeurl,
		&info.Message,
		&info.Status,
	)
	info.Jid = 1

	// if err != nil {
	// 	return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	// }

	if err != nil || info.Status >= 2 {
		return response.Resp().Json(gin.H{"msg": "沒有此活動"}, http.StatusNotFound)
	}

	return response.Resp().Json(gin.H{"msg": info}, http.StatusOK)
}
