package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// type SendMsg struct {
// 	UID       string `binding:"required"`
// 	PWD       string `binding:"required"`
// 	SB        string `json:",omitempty"`
// 	MSG       string `binding:"required"`
// 	DEST      string `binding:"required"`
// 	ST        string `json:",omitempty"`
// 	RETRYTIME string `json:",omitempty"`
// }

//聊天室簡訊發送
func SendSMS(context *gin.Context) *response.Response {
	token := context.GetString("token")
	if !CheckConnection(token) {
		return response.Resp().Json(gin.H{"msg": "token無效,請重新登入"}, http.StatusUnauthorized)
	}

	chatroomId := context.Param("chatroomId")
	text := context.DefaultPostForm("text", "")
	Log.Infof("[SendSMS] chatroomId:%s | text:%s", chatroomId, text)

	if IsEmpty(chatroomId) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	var (
		mobile  int64
		eventID int
	)
	qs := "SELECT eventID, mobile FROM chatroom WHERE chatroomID = ?"
	err := DB.Raw(qs, chatroomId).Row().Scan(
		&eventID,
		&mobile,
	)

	if err != nil {
		return response.Resp().Json(gin.H{"msg": "聊天室不存在"}, http.StatusNotFound)
	}

	//get code
	chatToken := fmt.Sprintf("%s%s", fmt.Sprint(eventID), chatroomId)
	var shortCode string
	code, err := getCode(chatToken)
	if err != nil {
		Log.Errorf("[getCode] %s", err)
	}
	if code != nil {
		shortCode = string(code)
	} else {
		shortCode = chatToken
	}

	chatroomUrl := fmt.Sprintf("%s%s", Config.GetString("chatroomUrl"), shortCode)
	msg := fmt.Sprintf("%s\n%s", text, chatroomUrl)
	dest := fmt.Sprint(mobile)
	url := fmt.Sprintf("https://%s/API21/HTTP/SendSMS.ashx", Config.GetString("smsAPI.ip"))

	data := make(map[string]string)
	data["MSG"] = msg
	data["DEST"] = dest
	Log.Infof("[SendSMS] data:%s", data)

	postData := new(bytes.Buffer)
	w := multipart.NewWriter(postData)
	for k, v := range data {
		w.WriteField(k, v)
	}
	w.Close()

	//簡訊發送
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
		Log.Errorf("[SendSMS] error: resp is %s", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Printf("%s\n", body)
	Log.Infof("[SendSMS] result:%s", body)

	result := strings.Split(string(body), ",")
	//檢查是否有資料
	param := StrToInt(result[0])
	if param < 0 {
		return response.Resp().Json(gin.H{"msg": string(body)}, http.StatusOK)
	}

	//add smsmap
	batch := result[4]
	AddSMS(batch, context.GetInt("id"), eventID, chatroomId, msg, "", StrToInt(SMS), 1)

	point := StrToFloat(result[0])

	return response.Resp().Json(gin.H{"point": point}, http.StatusOK)
}

func getCode(token string) (code []byte, err error) {
	influx := GetInflux()
	defer influx.Close()
	qs := fmt.Sprintf("SELECT shortCode, mobile FROM tokenmap WHERE chatToken = '%s'", token)

	resp, err := QueryInflux(influx, qs)
	if err != nil {
		return nil, err
	}

	if len(resp[0].Series) == 0 {
		return nil, nil
	}

	code = []byte(fmt.Sprint(resp[0].Series[0].Values[0][1]))

	return code, nil
}
