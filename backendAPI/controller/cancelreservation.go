package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//取消預約
func CancelReservation(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	token := context.GetString("token")

	if !CheckConnection(token) {
		return response.Resp().Json(gin.H{"msg": "token無效,請重新登入"}, http.StatusUnauthorized)
	}

	bid := context.PostForm("bid")
	mtype := context.PostForm("type")

	Log.Infof("[CancelReservation] bid:%s", bid)

	if IsEmpty(bid) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	var url string
	if mtype == SMS {
		url = fmt.Sprintf("https://%s/API21/HTTP/EraseBooking.ashx", Config.GetString("smsAPI.ip"))
	} else if mtype == MMS {
		url = fmt.Sprintf("https://%s/API21/HTTP/MMS/EraseBooking.ashx", Config.GetString("smsAPI.ip"))
	} else {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	data := make(map[string]string)
	data["BID"] = bid

	body, err := SendRequest(data, url, token)
	if err != nil {
		Log.Errorf("[CancelReservation] %s", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	Log.Infof("[CancelReservation] result:%s", body)

	result := strings.Split(string(body), ",")

	//檢查是否成功
	param := StrToInt(result[0])
	var point float64
	if param >= 0 {
		point = StrToFloat(result[1])
	} else {
		return response.Resp().Json(gin.H{"msg": string(body)}, http.StatusBadRequest)
	}
	return response.Resp().Json(gin.H{"point": point}, http.StatusOK)
}
