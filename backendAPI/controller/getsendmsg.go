package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	success           = 0 //傳送成功
	sending           = 1 //傳送中
	reservation       = 2 //預約
	cancelReservation = 3 //取消預約
	fail              = 4 //傳送失敗
	formatError       = 5 //格式錯誤
	emptyNumber       = 6 //空號
	timeout           = 7 //逾時
)

type sendMsg struct {
	Time    int64  `json:"time"`
	Bid     string `json:"bid"`
	Text    string `json:"text"`
	Subject string `json:"subject"`
	Count   int    `json:"count"`
	Status  int    `json:"status"`
	Cost    int    `json:"cost"`
	Mobile  string `json:"mobile"`
}

type respData struct {
	SMS_COUNT int
	BID       string
	DATA      []map[string]interface{}
}

//簡訊發送紀錄查詢
func GetSendMsg(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	token := context.GetString("token")
	if !CheckConnection(token) {
		return response.Resp().Json(gin.H{"msg": "token無效,請重新登入"}, http.StatusUnauthorized)
	}

	mtype := context.PostForm("type")
	start := int64(StrToInt(context.DefaultPostForm("start", "0")))
	end := int64(StrToInt(context.DefaultPostForm("end", "0")))

	Log.Infof("[GetSendMsg] mtype:%s", mtype)

	//條件查詢
	is_mobile := false
	is_status := false

	mobile := context.PostForm("mobile")
	statuslist := context.PostForm("status")
	status := -1
	if !IsEmpty(statuslist) {
		status = StrToInt(statuslist)
	}

	if !IsEmpty(mobile) {
		is_mobile = true
		if !IsMobile(mobile) {
			return response.Resp().Json(gin.H{"msg": "手機號碼格式錯誤"}, http.StatusBadRequest)
		}
	}

	if status >= 0 {
		is_status = true
		if status > 7 {
			return response.Resp().Json(gin.H{"msg": "狀態參數錯誤"}, http.StatusBadRequest)
		}
	}

	//檢查時間參數
	if start != 0 && end != 0 && start < end {
		if !isTimeRange(start) {
			return response.Resp().Json(gin.H{"msg": "時間參數錯誤"}, http.StatusBadRequest)
		}
	} else if start == 0 && end == 0 {
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).UnixNano()
		end = now.UnixNano()
	} else {
		return response.Resp().Json(gin.H{"msg": "時間參數錯誤"}, http.StatusBadRequest)
	}

	var url string
	data := make(map[string]string)

	//send getDR API
	if mtype == SMS {
		url = fmt.Sprintf("https://%s/API21/HTTP/GetDeliveryStatus.ashx", Config.GetString("smsAPI.ip"))
	} else if mtype == MMS {
		url = fmt.Sprintf("https://%s/API21/HTTP/MMS/GetDeliveryStatus.ashx", Config.GetString("smsAPI.ip"))
	} else {
		return response.Resp().Json(gin.H{"msg": "參數錯誤"}, http.StatusBadRequest)
	}
	// fmt.Printf("url:%s\n", url)

	data["PNO"] = "1"
	data["RESPFORMAT"] = "1"

	//search time range (batch,text)
	batchs, err := SearchBatchs(StrToInt(mtype), context.GetInt("id"), start, end)
	if err != nil {
		Log.Errorf("[SearchBatchs] %s\n", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	var s sendMsg
	var list []sendMsg

	done := make(chan struct{})
	// 新增阻塞chan
	// errChan := make(chan error)

	pool := NewGoroutinePool(1)

	for _, b := range batchs {
		pool.Add(1)

		data["BID"] = b.Batch
		text := b.Text
		subject := b.Subject

		Log.Infof("[GetSendMsg] data:%s", data)

		go func() {
			body, err := SendRequest(data, url, token)
			if err != nil {
				Log.Errorf("[GetSendMsg] %s", err)
			} else {
				Log.Infof("[GetSendMsg] result:%s", json.RawMessage(body))
				var respData respData
				json.Unmarshal(body, &respData)

				// fmt.Println("respData-> ", respData)
				if respData.SMS_COUNT > 0 {
					s.Count = respData.SMS_COUNT
					s.Bid = respData.BID

					data := respData.DATA
					s.Time = TimeToUnix("2006-01-02 15:04:05", data[0]["SEND_TIME"].(string), "-8h")
					s.Subject = subject
					s.Text = text
					s.Cost = int(StrToFloat(data[0]["COST"].(string)))
					s.Mobile = data[0]["MOBILE"].(string)
					s.Status = sortStatus(StrToInt(data[0]["STATUS"].(string)))

					exists := true
					//時間查詢條件
					if is_mobile && !filterByMobile(mobile, s.Mobile) {
						exists = false
					}

					//狀態查詢條件
					if exists && is_status && !filterByStatus(status, s.Status) {
						exists = false
					}

					if exists {
						list = append(list, s)
					}
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
		return response.Resp().Json(gin.H{"list": list}, http.StatusOK)
	}
}

//檢查時間範圍
func isTimeRange(startTime int64) bool {
	//檢查是否在三個月內
	now := time.Now()
	timeRange := now.AddDate(0, -3, 0).UnixNano()
	if startTime < timeRange {
		return false
	}
	return true
}

type smsInfo struct {
	Batch   string `json:"batch"`
	Text    string `json:"text"`
	Subject string `json:"subject"`
}

func SearchBatchs(mtype, accountId int, startTime, endTime int64) ([]smsInfo, error) {
	influx := GetInflux()
	defer influx.Close()

	qs := fmt.Sprintf("SELECT batch, text, subject FROM smsmap WHERE type = %d AND accountID = %d AND time > %d AND time < %d ORDER BY DESC", mtype, accountId, startTime, endTime)

	resp, err := QueryInflux(influx, qs)
	if err != nil {
		return nil, err
	}

	fmt.Println("resp", resp)

	var s smsInfo
	var list []smsInfo

	if len(resp[0].Series) == 0 {
		return list, nil
	}

	values := resp[0].Series[0].Values
	for _, row := range values {
		s.Batch = fmt.Sprint(row[1])
		s.Text = strings.Replace(fmt.Sprint(row[2]), "\n", " ", -1)
		var sb string
		if strings.Compare(fmt.Sprint(row[3]), "<nil>") != 0 {
			sb = fmt.Sprint(row[3])
		}
		s.Subject = sb
		list = append(list, s)
	}

	return list, nil
}

func filterByMobile(filter, mobile string) bool {
	if strings.Compare(filter[len(filter)-9:], mobile[len(mobile)-9:]) == 0 {
		return true
	}
	return false
}

func filterByStatus(filter, status int) bool {
	if filter == status {
		return true
	}
	return false
}

func sortStatus(status int) int {
	var s int
	switch status {
	case 100, 700:
		s = success
	case 0:
		s = sending
	case 300:
		s = reservation
	case 303:
		s = cancelReservation
	case -1, -2, -5, 101, 102, 104, 105, 106, 107, 500:
		s = fail
	case -3:
		s = formatError
	case 103:
		s = emptyNumber
	case -4:
		s = timeout
	default:
		s = -1
	}
	return s
}
