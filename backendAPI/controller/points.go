package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

//點數扣除
func Points(context *gin.Context) *response.Response {
	token := context.GetString("token")
	if !CheckConnection(token) {
		return response.Resp().Json(gin.H{"msg": "token無效,請重新登入"}, http.StatusUnauthorized)
	}

	bid := context.PostForm("bid")
	userNo := context.PostForm("userNo")
	mr := context.PostForm("mr")
	deductTime := context.PostForm("deductTime")
	points := context.PostForm("points")

	if IsEmpty(bid) || IsEmpty(userNo) || IsEmpty(mr) || IsEmpty(deductTime) || IsEmpty(points) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	info := make(map[string]interface{})
	info["HandlerType"] = 2
	info["MaintType"] = 1
	info["CreditType"] = 5
	info["TransType"] = "091"
	info["UserNo"] = StrToInt(userNo)
	info["BID"] = bid
	info["MR"] = mr
	info["DeductTime"] = deductTime
	f := StrToFloat(points)
	decimal.DivisionPrecision = 2
	info["Points"] = decimal.NewFromFloat(f)

	timeFormat := fmt.Sprintf("%s-%s-%s %s:%s:%s", deductTime[:4], deductTime[4:6], deductTime[6:8], deductTime[8:10], deductTime[10:12], deductTime[12:14])
	info["Description"] = fmt.Sprintf("客戶於%s點擊聊天室 (批號:%s；識別號:%s)", timeFormat, bid, mr)

	jsonBytes, _ := json.Marshal(info)
	fmt.Printf("%s\n", jsonBytes)

	Log.Infof("[Points] jsonBytes:%s", jsonBytes)

	url := fmt.Sprintf("https://%s/API21/HTTP/PointsHandler.ashx", Config.GetString("smsAPI.ip"))

	var reqData = []byte(jsonBytes)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

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
		Log.Errorf("[Points] error: resp is %s", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Printf("%s\n", body)

	// Log.Infof("[Points] result:%s", body)

	respData := make(map[string]interface{})
	json.Unmarshal([]byte(body), &respData)

	if !respData["Result"].(bool) {
		return response.Resp().Json(gin.H{"msg": respData["Msg"]}, http.StatusNotFound)
	}
	// result := fmt.Sprint(respData["Result"])
	// if strings.Compare(result, "false") == 0 {
	// 	return response.Resp().Json(gin.H{"msg": respData["Msg"]}, http.StatusNotFound)
	// }
	return response.Resp().Json(gin.H{"msg": respData["Msg"]}, http.StatusOK)
}
