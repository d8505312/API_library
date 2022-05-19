package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

var (
	SMS = "0"
	MMS = "1"
)

//簡訊批量發送
func SendBatchSMS(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	token := context.GetString("token")
	if !CheckConnection(token) {
		return response.Resp().Json(gin.H{"msg": "token無效,請重新登入"}, http.StatusUnauthorized)
	}

	eventId := context.Param("eventId")
	text := context.PostForm("text")
	mtype := context.PostForm("type")
	subject := context.DefaultPostForm("subject", "")
	list := context.PostForm("list")
	sendTime := StrToInt(context.DefaultPostForm("sendTime", "0"))
	phrases := context.DefaultPostForm("phrases", "")

	Log.Infof("[SendBatchSMS] eventId:%s | text:%s | mtype:%s | subject:%s | list:%s | sendTime:%s", eventId, text, mtype, subject, list, sendTime)

	if IsEmpty(text) || IsEmpty(list) {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	var url string
	dest := fmt.Sprint(list)

	dests := strings.Split(string(dest), ",")

	var data = make(map[string]string)
	//SMS不帶主旨
	if mtype == SMS {
		url = fmt.Sprintf("https://%s/API21/HTTP/SendSMS.ashx", Config.GetString("smsAPI.ip"))
		if !IsEmpty(subject) {
			data["SB"] = subject
		}
		//MMS必帶主旨
	} else if mtype == MMS {
		if IsEmpty(subject) || !checkSubject(subject) {
			return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
		}
		file, err := context.FormFile("image")
		if err != nil {
			Log.Error("[MMS] ", err)
			return response.Resp().Json(gin.H{"msg": "檔案不存在"}, http.StatusNotFound)
		}
		fileTmp, _ := file.Open()
		byteFile, _ := ioutil.ReadAll(fileTmp)

		url = fmt.Sprintf("https://%s/API21/HTTP/MMS/SendMMS.ashx", Config.GetString("smsAPI.ip"))

		data["SB"] = subject
		data["ATTACHMENT"] = base64.StdEncoding.EncodeToString(byteFile)
		data["TYPE"] = strings.Split(file.Filename, ".")[1]
	} else {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}

	if sendTime != 0 {
		data["ST"] = UnixToStr(int64(sendTime), "20060102150405")
	}

	point := -1.0
	done := make(chan struct{})
	// 新增阻塞chan
	// errChan := make(chan error)

	pool := NewGoroutinePool(1)

	for _, d := range dests {
		pool.Add(1)

		data["DEST"] = fmt.Sprintf("886%s", d[len(d)-9:])

		var addtoken bool
		var shortCode string
		//檢查電話號是否使用過
		code, err := checkMobileUse(eventId, data["DEST"])
		if err != nil {
			Log.Errorf("[checkMobileUse] %s", err)
		}
		if code != nil {
			shortCode = string(code)
			addtoken = false
		} else {
			shortCode = GetShortCode(6)
			addtoken = true
		}

		chatroomId := GetTokenString(7)
		chatToken := fmt.Sprintf("%s%s", eventId, chatroomId)
		chatroomUrl := fmt.Sprintf("%s%s", Config.GetString("chatroomUrl"), shortCode)
		var msg string
		if IsEmpty(phrases) {
			msg = fmt.Sprintf("%s\n%s", text, chatroomUrl)
		} else {
			msg = fmt.Sprintf("%s\n%s %s", text, phrases, chatroomUrl)
		}

		data["MSG"] = msg

		Log.Infof("[SendBatchSMS] data:%s", data)

		go func() {
			body, err := SendRequest(data, url, token)
			if err != nil {
				Log.Errorf("[SendBatchSMS] %s", err)
			} else {
				Log.Infof("[SendBatchSMS] result:%s", body)

				result := strings.Split(string(body), ",")

				//檢查是否有資料
				// if len(result) <= 2 {
				param := StrToFloat(result[0])
				if param >= 0 {
					point = StrToFloat(result[0])

					if addtoken {
						//add tokenmap
						AddChatToken(shortCode, chatToken, data["DEST"])
					}

					//add smsmap
					batch := result[4]
					AddSMS(batch, context.GetInt("id"), StrToInt(eventId), chatroomId, msg, subject, StrToInt(mtype), 1)
				} else {
					//簡訊發送失敗
					//todo 錯誤處理
				}
			}

			pool.Done()
			// if err != nil {
			// 	errChan <- errors.New("error")
			// }
		}()
	}

	go func() {
		pool.Wait()
		close(done)
	}()

	select {
	// 錯誤快返回
	// case err := <-errChan:
	// 	return response.Resp().Json(gin.H{"msg": err}, http.StatusInternalServerError)
	case <-time.After(30 * time.Second):
		return response.Resp().Json(gin.H{"msg": "請求超時"}, http.StatusRequestTimeout)
	case <-done:
		if point < 0 {
			return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
		}
		return response.Resp().Json(gin.H{"point": point}, http.StatusOK)
	}
}

//檢查此門號是否有此活動連接
func checkMobileUse(eventid, mobile string) (code []byte, err error) {
	influx := GetInflux()
	defer influx.Close()
	phone := StrToInt(mobile)
	qs := fmt.Sprintf("SELECT chatToken, shortCode, mobile FROM tokenmap WHERE mobile= %d", phone)

	resp, err := QueryInflux(influx, qs)
	if err != nil {
		return nil, err
	}

	if len(resp[0].Series) == 0 {
		return nil, nil
	}

	values := resp[0].Series[0].Values
	for _, row := range values {
		tokn := fmt.Sprint(row[1])
		evnid := tokn[0 : len(tokn)-7]
		if strings.Compare(evnid, eventid) == 0 {
			code = []byte(fmt.Sprint(row[2]))
			break
		}
	}
	return code, nil
}

//檢查主旨是否符合限制
func checkSubject(s string) bool {
	//主旨長度不超過39字
	if utf8.RuneCountInString(s) > 39 {
		return false
	}
	//主旨有中文,長度不超過15字
	if IsChinese(s) && utf8.RuneCountInString(s) > 15 {
		return false
	}
	return true
}

//新增ChatToken
func AddChatToken(code, token, mobile string) {
	influx := GetInflux()
	defer influx.Close()

	tags := map[string]string{
		"chatToken": token,
		"shortCode": code,
	}

	fields := map[string]interface{}{}
	if !IsEmpty(mobile) {
		fields["mobile"] = StrToInt(mobile)
	}

	insErr := InsertInflux(influx, "tokenmap", tags, fields)
	if insErr != nil {
		Log.Errorf("insErr:%s", insErr.Error())
		return
	}
	Log.Info("Insert ChatToken Ok")
}

//新增SMS
func AddSMS(batch string, accountID, eventId int, chatroomId, text, subject string, mtype, count int) {
	influx := GetInflux()
	defer influx.Close()

	tags := map[string]string{
		"batch": batch,
	}

	fields := map[string]interface{}{}
	fields["accountID"] = accountID
	fields["eventID"] = eventId
	fields["chatroomID"] = chatroomId
	fields["type"] = mtype
	fields["count"] = count
	fields["text"] = text
	if !IsEmpty(subject) {
		fields["subject"] = subject
	}

	insErr := InsertInflux(influx, "smsmap", tags, fields)
	if insErr != nil {
		Log.Errorf("insErr:%s", insErr.Error())
		return
	}
	Log.Info("Insert SMS Ok")
}

//新增ChatRoom
// func addChatRoom(chatroomID, eventID, mobile string) {
// 	//insert chatroom table
// 	fmt.Printf("eventID:%s mobile:%s\n", eventID, mobile)
// 	i_qs := "INSERT INTO chatroom(chatroomID, eventID, mobile, lastVisit, lastVisitB) VALUES(?,?,?,?,?)"
// 	ins := DB.Exec(i_qs, &chatroomID, &eventID, &mobile, time.Now().UTC(), time.Now().UTC())
// 	if ins.Error != nil {
// 		Log.Errorf("insErr:%s", ins.Error)
// 		return
// 	}
// 	Log.Info("Insert  ChatRoom Ok")
// }
