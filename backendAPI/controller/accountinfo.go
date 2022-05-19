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

//簡訊帳號查詢
func AccountInfo(token string) (respData []map[string]interface{}, err error) {

	info := make(map[string]interface{})
	info["HandlerType"] = 1
	info["GetInfoType"] = 3

	jsonBytes, _ := json.Marshal(info)
	fmt.Printf("%s\n", jsonBytes)

	url := fmt.Sprintf("https://%s/API21/HTTP/AccountHandler.ashx", Config.GetString("smsAPI.ip"))

	var reqData = []byte(jsonBytes)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, fmt.Errorf("error: req is %s", err)
	}
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
		return nil, fmt.Errorf("error: resp is %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error: status is not ok,resp code is %d, body is %s", resp.StatusCode, string(body))
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	Log.Infof("[AccountInfo] result:%s", json.RawMessage(body))

	json.Unmarshal([]byte(body), &respData)

	return respData, nil
}

func AccountValueExist(maps []map[string]interface{}, key string, value int) bool {
	for _, v := range maps {
		if v[key] == value {
			return true
		}
	}
	return false
}

func GetRoleInfo(maps []map[string]interface{}, userNo int) map[string]interface{} {
	for _, v := range maps {
		if int(v["UserNo"].(float64)) == userNo {
			return v
		}
	}
	return nil
}
