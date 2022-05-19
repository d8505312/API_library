package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//簡訊餘額查詢
func GetCredit(context *gin.Context) *response.Response {
	token := context.GetString("token")

	if !CheckConnection(token) {
		return response.Resp().Json(gin.H{"msg": "token無效,請重新登入"}, http.StatusUnauthorized)
	}

	url := fmt.Sprintf("https://%s/API21/HTTP/GetCredit.ashx", Config.GetString("smsAPI.ip"))

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "text/plain")

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
		Log.Errorf("[GetCredit] error: resp is %s", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Printf("%s\n", body)

	Log.Infof("[GetCredit] result:%s", body)

	result := strings.Split(string(body), ",")
	//檢查是否有資料
	if len(result) > 1 {
		return response.Resp().Json(gin.H{"status": result[0], "msg": result[1]}, http.StatusOK)
	}

	point := StrToFloat(result[0])

	return response.Resp().Json(gin.H{"point": point}, http.StatusOK)
}
