package controller

import (
	"backendAPI/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(context *gin.Context) *response.Response {
	token := context.GetString("token")
	if CheckConnection(token) {
		CloseConnection(token)
	}
	return response.Resp().Json(gin.H{"msg": "OK"}, http.StatusOK)
}
