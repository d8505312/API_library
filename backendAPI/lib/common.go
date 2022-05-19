package lib

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

func IsEvent(admin, accountid, eventid int) bool {
	if admin == Manager {
		return true
	}
	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM accountEvent WHERE eventID = ? AND accountID = ?)"
	if err := DB.Raw(s_qs, eventid, accountid).Row().Scan(&exists); err != nil {
		Log.Error("[IsEvent] ", err)
		return false
	}
	return exists
}

func IsChatRoom(admin, accountid int, chatroomid string) bool {
	if admin == Manager {
		return true
	}
	var exists bool
	s_qs := "SELECT EXISTS(SELECT accountID FROM accountEvent WHERE eventID = (SELECT eventID FROM chatroom WHERE chatroomID = ?) AND accountID = ?)"
	if err := DB.Raw(s_qs, chatroomid, accountid).Row().Scan(&exists); err != nil {
		Log.Error("[IsChatRoom] ", err)
		return false
	}
	return exists
}

//簡訊相關API發送
func SendRequest(data map[string]string, url, token string) (body []byte, err error) {
	// fmt.Printf("sendRequest->%v\n", data)

	postData := new(bytes.Buffer)
	w := multipart.NewWriter(postData)
	for k, v := range data {
		w.WriteField(k, v)
	}
	w.Close()

	req, err := http.NewRequest("POST", url, postData)
	if err != nil {
		return nil, fmt.Errorf("error: req is %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", w.FormDataContentType())

	fmt.Printf("url:%s\n", url)

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
}

func AddAccount(customerId, admin int, account, name string) (accountid int, err error) {
	_, publicKey, err := KeyPairs()
	if err != nil {
		return 0, fmt.Errorf("[AddAccount] GenKey Failed!  %s\n", err)
	}

	keyExpire := (int)(time.Now().Unix() + Config.GetInt64("tokenExpirePeriod"))

	tx := DB.Begin()

	// insert account table
	i_qs := "INSERT INTO account(customerID, account, name, jwtKey, keyExpire,admin) VALUES(?,?,?,?,?,?)"
	ins := tx.Exec(i_qs, &customerId, &account, &name, &publicKey, &keyExpire, &admin)
	if ins.Error != nil {
		tx.Rollback()
		return 0, fmt.Errorf("[AddAccount] DB INSERT Failed!  %s\n", ins.Error)
	}

	// get last accountID
	var id uint64
	err = tx.Raw("SELECT LAST_INSERT_ID()").Row().Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("[AddAccount] Get Last Id Failed!  %s\n", err)
	}

	accountid = int(id)
	tx.Commit()

	Log.Infof("[AddAccount] INSERT OK!\n")

	return accountid, nil
}
