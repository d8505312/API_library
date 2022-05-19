package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*	Report 參數說明

	BatchID	批次代碼(發送代碼)
	RM		手機號碼
	RT		回報時間，格式為yyyyMMddHHmmss
	STATUS	DR狀態
	SM		回覆之簡訊內容
	MR		Record No.每一筆簡訊之發送代碼
	ST		發送時間
	SUBJECT	簡訊主旨
	NAME	姓名
	SOURCE	發送來源
	CHARGE	計點點數
	TYPE	SMS、MMS
	MSGID	發送時所填入的MSGID
*/

// DR/MO 回報
func Report(context *gin.Context) *response.Response {
	batchId := context.Query("BatchID")
	rm := context.Query("RM")
	rt := context.Query("RT")
	status := context.Query("STATUS")
	sm := context.Query("SM")
	mr := context.Query("MR")
	st := context.Query("ST")
	subject := context.Query("SUBJECT")
	name := context.Query("NAME")
	source := context.Query("SOURCE")
	charge := context.Query("CHARGE")
	mtype := context.Query("TYPE")
	msgId := context.Query("MSGID")

	fmt.Printf("batchId:%s rm:%s rt:%s status:%s sm:%s mr:%s st:%s subject:%s name:%s source:%s charge:%s mtype:%s msgId:%s", batchId, rm, rt, status, sm, mr, st, subject, name, source, charge, mtype, msgId)

	// 999為回覆簡訊
	if StrToInt(status) == 999 {

		//去除空白及換行符
		sm = strings.Replace(sm, " ", "", -1)
		sm = strings.Replace(sm, "\n", "", -1)
		smLen := len(sm)
		first := sm[:1]

		//檢查開頭是否為1, 1為mo簡訊驗證代碼
		if StrToInt(first) == 1 && smLen == 6 {
			rm = strings.TrimSpace(rm)
			//驗證手機號碼格式
			if IsMobile(rm) {
				rm = fmt.Sprintf("886%s", rm[len(rm)-9:])
				verifyCode(StrToInt(sm), StrToInt(rm))
			}
		}
	}

	return response.Resp().Json(gin.H{"Msg": "OK"}, http.StatusOK)
}

type Device struct {
	Key  string `json:"k"`
	Time int    `json:"t"`
	Code int    `json:"c,omitempty"`
}

//驗證簡訊驗證碼
func verifyCode(code, mobile int) {
	fmt.Println("code", code, "mobile", mobile)

	tx := DB.Begin()

	var uid int
	var device json.RawMessage
	s_qs := "SELECT uid, userDevice FROM user WHERE mobile = ? LIMIT 1"
	if err := tx.Raw(s_qs, mobile).Row().Scan(&uid, &device); err != nil {
		tx.Rollback()
		Log.Error("[verifyCode] ", err)
		return
	}

	//解析device
	var devices []Device
	json.Unmarshal(device, &devices)

	var devs []Device
	var new_devs []Device

	for _, d := range devices {
		if d.Code > 0 {
			new_devs = append(new_devs, d)
		} else {
			devs = append(devs, d)
		}
	}

	if len(new_devs) == 0 {
		Log.Warnf("[verifyCode] 尚無新裝置")
		return
	}

	//檢查待驗證設備驗證碼
	if code == new_devs[0].Code {
		//檢查驗證碼是否過期
		now := time.Now().UTC()
		deadline := 600
		exp := int64(new_devs[0].Time + deadline)
		if now.Unix() > exp {
			Log.Warnf("[verifyCode] 驗證碼過期")
			return
		}
	} else {
		Log.Warnf("[verifyCode] 驗證碼錯誤")
		return
	}

	new_dev := Device{
		Key:  new_devs[0].Key,
		Time: new_devs[0].Time,
	}
	//排序
	sort.Slice(devs, func(i, j int) bool {
		return devs[i].Time < devs[j].Time
	})

	//未滿三組則新增, 否則取代最久未使用之裝置
	if len(devs) < 3 {
		devs = append(devs, new_dev)
	} else {
		devs[0] = new_dev
	}
	encode, _ := json.Marshal(devs)
	u_qs := "UPDATE user SET userDevice = ? WHERE uid = ? AND mobile = ?"
	ins := tx.Exec(u_qs, encode, uid, mobile)
	if ins.Error != nil {
		tx.Rollback()
		Log.Errorf("[verifyCode] ins err:%s", ins.Error)
		return
	}

	tx.Commit()

	Log.Info("[verifyCode] Update user Ok")
}
