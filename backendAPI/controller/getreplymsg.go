package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//簡訊回覆紀錄查詢
func GetReplyMsg(context *gin.Context) *response.Response {
	token := context.GetString("token")
	if !CheckConnection(token) {
		return response.Resp().Json(gin.H{"msg": "token無效,請重新登入"}, http.StatusUnauthorized)
	}

	bid := context.PostForm("bid")
	pno := context.DefaultPostForm("pno", "1")

	if IsEmpty(bid) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	data := make(map[string]string)

	data["BID"] = bid
	data["PNO"] = pno
	data["RESPFORMAT"] = "1"

	Log.Infof("[GetReplyMsg] data:%s", data)

	url := fmt.Sprintf("https://%s/API21/HTTP/GetReplyMessage.ashx", Config.GetString("smsAPI.ip"))

	fmt.Printf("url:%s\n", url)

	postData := new(bytes.Buffer)
	w := multipart.NewWriter(postData)
	for k, v := range data {
		w.WriteField(k, v)
	}
	w.Close()

	req, _ := http.NewRequest("POST", url, postData)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{
		Timeout: time.Second * time.Duration(Config.GetInt64(("smsAPI.delay"))),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: Config.GetBool("isSkipVerify"),
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		Log.Errorf("[GetReplyMsg] error: resp is %s", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println("resp Body:", string(body))

	var mapResult map[string]interface{}
	json.Unmarshal([]byte(body), &mapResult)

	return response.Resp().Json(mapResult, http.StatusOK)
}
