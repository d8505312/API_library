package controller

//controller/Index.go

import (
	"backendAPI/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(context *gin.Context) *response.Response {

	if 1+1 == 2 {

		return response.Resp().Json(gin.H{"msg": "hello world"}, http.StatusOK)

	}

	return response.Resp().Json(gin.H{"msg": "hello world"}, http.StatusOK)
}
