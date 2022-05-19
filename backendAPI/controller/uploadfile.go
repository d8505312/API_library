package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//檔案上傳
func UpLoadFile(context *gin.Context) *response.Response {

	file, err := context.FormFile("file")
	if err != nil {
		return response.Resp().Json(gin.H{"msg": "檔案不存在"}, http.StatusNotFound)
	}

	id := GetFileID()
	fileExt := GetFileExt(file.Filename)
	fileName := ChangeFileName(file.Filename, id)
	fileId, _ := strconv.Atoi(id)
	fmt.Println("FileName: ", fileName)
	Log.Info("[UpLoadFile] FileName:", fileName)

	if err := context.SaveUploadedFile(file, Config.GetString("fileDir")+fileName); err != nil {
		Log.Error("[UpLoadFile] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	//設定有效期限
	now := time.Now()
	deadline, _ := time.ParseDuration("72h")
	exp := now.Add(deadline).Unix()

	return response.Resp().Json(gin.H{"fileid": fileId, "ext": fileExt, "exp": exp}, http.StatusOK)
}
