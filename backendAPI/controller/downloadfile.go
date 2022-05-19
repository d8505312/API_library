package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ImageFileData struct {
	Fileid      int
	ShowName    string
	Description string
	Width       int
	Height      int
	FileSize    int
}

//檔案下載
func DownLoadFile(context *gin.Context) *response.Response {

	fileId := context.Param("fileid")
	Log.Infof("[DownLoadFile] fileId:%s", fileId)

	influx := GetInflux()
	defer influx.Close()
	fs := fmt.Sprintf(":%s,", fileId)
	qs := fmt.Sprintf("SELECT format FROM message WHERE msgType>=2 AND format =~ /%s/", fs)

	resp, err := QueryInflux(influx, qs)
	if err != nil {
		Log.Error("[DownLoadFile] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	if len(resp[0].Series) == 0 {
		return response.Resp().Json(gin.H{"msg": "檔案不存在"}, http.StatusNotFound)
	}

	format := fmt.Sprint(resp[0].Series[0].Values[0][1])

	var f ImageFileData
	json.NewDecoder(strings.NewReader(format)).Decode(&f)
	fmt.Println("format", f.ShowName)
	Log.Infof("[DownLoadFile] ShowName:%s", f.ShowName)

	fileDir := Config.GetString("fileDir")
	fileShowName := f.ShowName
	fileName := ChangeFileName(fileShowName, fileId)
	filePath := fileDir + fileName
	fileTmp, errByOpenFile := os.Open(filePath)

	if IsEmpty(fileDir) || IsEmpty(fileName) || errByOpenFile != nil {
		return response.Resp().Json(gin.H{"msg": "檔案不存在"}, http.StatusNotFound)
	}

	defer fileTmp.Close()

	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Type", "application/x-download")
	context.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileShowName))
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Cache-Control", "no-cache")
	context.File(filePath)

	Log.Info("download OK")

	return response.Resp()
}
