package controller

import (
	. "backendAPI/lib"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// var SMSUrl = fmt.Sprintf("https://%s/API21/HTTP/ConnectionHandler.ashx", Config.GetString("smsAPI.ip"))

//取得簡訊平台服務連線
func GetConnection(account, password string) (body []byte, err error) {
	uid := account
	pwd := password

	info := make(map[string]interface{})
	info["HandlerType"] = 3
	info["VerifyType"] = 1
	info["UID"] = uid
	info["PWD"] = pwd

	jsonBytes, _ := json.Marshal(info)
	fmt.Printf("%s\n", jsonBytes)
	Log.Infof("[GetConnection] jsonBytes:%s", jsonBytes)

	url := fmt.Sprintf("https://%s/API21/HTTP/ConnectionHandler.ashx", Config.GetString("smsAPI.ip"))

	var reqData = []byte(jsonBytes)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, fmt.Errorf("error: req is %s", err)
	}
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
		return nil, fmt.Errorf("error: resp is %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error: status is not ok,resp code is %d, body is %s", resp.StatusCode, string(body))
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return body, nil

	// result := fmt.Sprint(respData["Result"])
	// if strings.Compare(result, "false") == 0 {
	// 	return response.Resp().Json(gin.H{"msg": respData["Msg"]}, http.StatusNotFound)
	// }
}

//檢查連線狀態
func CheckConnection(token string) bool {

	info := make(map[string]interface{})
	info["HandlerType"] = 3
	info["VerifyType"] = 2

	jsonBytes, _ := json.Marshal(info)
	fmt.Printf("%s\n", jsonBytes)
	Log.Infof("[CheckConnection] jsonBytes:%s", jsonBytes)

	url := fmt.Sprintf("https://%s/API21/HTTP/ConnectionHandler.ashx", Config.GetString("smsAPI.ip"))

	var reqData = []byte(jsonBytes)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		Log.Errorf("[CheckConnection] error: resp is %s", err)
		return false
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	// fmt.Println(resp.StatusCode)
	// fmt.Printf("%s\n", body)

	Log.Infof("[CheckConnection] result:%s", json.RawMessage(body))

	respData := make(map[string]interface{})
	json.Unmarshal([]byte(body), &respData)
	result := respData["Result"].(bool)

	return result
}

//關閉簡訊平台服務連線
func CloseConnection(token string) bool {

	info := make(map[string]interface{})
	info["HandlerType"] = 3
	info["VerifyType"] = 3

	jsonBytes, _ := json.Marshal(info)
	// fmt.Printf("%s\n", jsonBytes)
	Log.Infof("[CloseConnection] jsonBytes:%s", jsonBytes)

	url := fmt.Sprintf("https://%s/API21/HTTP/ConnectionHandler.ashx", Config.GetString("smsAPI.ip"))

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
		Log.Errorf("[CloseConnection] error: resp is %s", err)
		return false
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	// fmt.Println(resp.StatusCode)
	// fmt.Printf("%s\n", body)

	Log.Infof("[CloseConnection] result:%s", json.RawMessage(body))

	respData := make(map[string]interface{})
	json.Unmarshal([]byte(body), &respData)
	result := respData["Result"].(bool)

	return result
}
