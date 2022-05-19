package controller

//controller/Token.go

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	. "backendAPI/lib"
	"backendAPI/response"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

var (
	Self = "0"
	E8D  = "1"
)

func Token(c *gin.Context) *response.Response {

	//var result Result
	// var id int
	var id uint64
	//var jwtKey string
	var privateKey string
	var publicKey string
	var keyExpire int
	var admin int
	var name string
	var ips string
	var token string

	//帳號來源
	domain := c.DefaultPostForm("domain", "1")

	account := c.PostForm("account")
	// pass := fmt.Sprintf("%x", sha256.Sum256([]byte(c.PostForm("password"))))
	pass := c.PostForm("password")
	// fmt.Printf("account:%s | pass:%s\n", account, pass)

	if IsEmpty(domain) || IsEmpty(account) || IsEmpty(pass) {
		return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
	}

	switch domain {
	case Self:
		token = ""
		sha256_pass := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
		err := DB.Raw("SELECT accountID,admin,name,ips FROM account WHERE account = ? and password=?", account, sha256_pass).Row().Scan(&id, &admin, &name, &ips)
		if err != nil {
			Log.Errorf("[Token] %s", err)
			return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
		}

		if len(ips) > 1 { //IP合法性檢查
			if strings.Contains(ips, c.ClientIP()) {
				return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusUnauthorized)
			}
		}

		//生成privateKey & publicKey
		privateKey, publicKey, err = KeyPairs()
		if err != nil {
			//fmt.Printf("GenKey Failed!  %s\n", err)
			Log.Errorf("[Token] %s", err)
			return response.Resp().Json(gin.H{"key": "", "expire": 97, "admin": 0}, http.StatusInternalServerError)
		}
		//fmt.Printf("%s\n%s\n", privateKey, publicKey)

		keyExpire = (int)(time.Now().Unix() + Config.GetInt64("tokenExpirePeriod"))

		tx := DB.Exec("update account set jwtKey=?, keyExpire=? where accountID=?", publicKey, keyExpire, id)

		if tx.Error != nil {
			//fmt.Printf("update Failed!  %s\n", err)
			return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusInternalServerError)
		}

	case E8D:
		tx := DB.Begin()

		// 檢查是否有帳號, 沒有則新增
		var exists bool
		s_qs := "SELECT EXISTS(SELECT account FROM account WHERE account = ?)"
		if err := tx.Raw(s_qs, account).Row().Scan(&exists); err != nil {
			tx.Rollback()
			Log.Errorf("[Token] %s", err)
			return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusInternalServerError)
		}

		//生成privateKey & publicKey
		var err error
		privateKey, publicKey, err = KeyPairs()
		if err != nil {
			//fmt.Printf("GenKey Failed!  %s\n", err)
			Log.Errorf("[Token] %s", err)
			return response.Resp().Json(gin.H{"key": "", "expire": 97, "admin": 0}, http.StatusInternalServerError)
		}
		//fmt.Printf("%s\n%s\n", privateKey, publicKey)

		keyExpire = (int)(time.Now().Unix() + Config.GetInt64("tokenExpirePeriod"))

		if exists {
			if err := tx.Raw("SELECT token FROM account WHERE account = ?", account).Row().Scan(&token); err != nil {
				tx.Rollback()
				Log.Errorf("[Token] %s", err)
				return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusInternalServerError)
			}

			//先關閉舊的e8d連線
			// if !IsEmpty(token) {
			// 	CloseConnection(token)
			// }

			//every8d 登入驗證
			body, err := GetConnection(account, pass)
			if err != nil {
				Log.Errorf("[GetConnection] %s", err)
				return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
			} else {
				Log.Infof("[GetConnection] result:%s", json.RawMessage(body))
			}
			respData := make(map[string]interface{})
			json.Unmarshal([]byte(body), &respData)

			if !respData["Result"].(bool) {
				return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
			}

			//every8d token
			token = fmt.Sprint(respData["Msg"])

			err = tx.Raw("SELECT accountID, admin, name, ips FROM account WHERE account = ?", &account).Row().Scan(&id, &admin, &name, &ips)
			if err != nil {
				// fmt.Printf(">>>>%s\n", err)
				tx.Rollback()
				Log.Errorf("[Token] %s", err)
				return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
			}

			if len(ips) > 1 { //IP合法性檢查
				if strings.Contains(ips, c.ClientIP()) {
					return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusUnauthorized)
				}
			}

			upd := tx.Exec("update account set jwtKey=?, keyExpire=?, token=? where accountID=?", &publicKey, &keyExpire, &token, &id)
			if upd.Error != nil {
				//fmt.Printf("update Failed!  %s\n", err)
				tx.Rollback()
				Log.Errorf("[Token] %s", err)
				return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusInternalServerError)
			}
		} else {
			//every8d 登入驗證
			body, err := GetConnection(account, pass)
			if err != nil {
				Log.Errorf("[GetConnection] %s", err)
				return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
			} else {
				Log.Infof("[GetConnection] result:%s", json.RawMessage(body))
			}
			resultData := make(map[string]interface{})
			json.Unmarshal([]byte(body), &resultData)

			if !resultData["Result"].(bool) {
				return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
			}

			//every8d token
			token = fmt.Sprint(resultData["Msg"])

			respData, err := AccountInfo(token)
			if err != nil {
				Log.Errorf("[AccountInfo] %s", err)
				return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
			}

			if len(respData) <= 0 {
				return response.Resp().Json(gin.H{"key": "", "expire": 98, "admin": 0}, http.StatusUnauthorized)
			}

			var customerAccount string
			var customerId uint64

			//檢查LoginName是否為自己, RoleID =1000 為母帳號 (admin=2)
			for _, d := range respData {
				if strings.Compare(d["LoginName"].(string), account) == 0 {
					name = d["UserName"].(string)
					if int(d["RoleID"].(float64)) == 1000 {
						admin = Manager
						break
					}
				} else {
					if int(d["RoleID"].(float64)) == 1000 {
						customerAccount = d["LoginName"].(string)
					}
				}
			}

			if admin == Manager {
				//insert customer table
				i_qs := "INSERT INTO customer(name, type, callable) VALUES(?,?,?)"
				ins := tx.Exec(i_qs, &account, &domain, 1)
				if ins.Error != nil {
					tx.Rollback()
					Log.Error("[Token] ", ins.Error)
					return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusInternalServerError)
				}

				// get last customerID
				err = tx.Raw("SELECT LAST_INSERT_ID()").Row().Scan(&customerId)
				if err != nil {
					tx.Rollback()
					Log.Error("[Token] ", err)
					return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusInternalServerError)
				}
			} else {
				err = tx.Raw("SELECT customerID FROM account WHERE account=?", customerAccount).Row().Scan(&customerId)
				if err != nil {
					tx.Rollback()
					Log.Error("[Token] ", err)
					return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusInternalServerError)
				}
			}

			// insert account table
			i_qs := "INSERT INTO account(customerID, account, name, jwtKey, keyExpire,token,admin) VALUES(?,?,?,?,?,?,?)"
			ins := tx.Exec(i_qs, &customerId, &account, &name, &publicKey, &keyExpire, &token, &admin)
			if ins.Error != nil {
				tx.Rollback()
				Log.Error("[Token] ", ins.Error)
				return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusInternalServerError)
			}

			// get last accountID
			err = tx.Raw("SELECT LAST_INSERT_ID()").Row().Scan(&id)
			if err != nil {
				tx.Rollback()
				Log.Error("[Token] ", err)
				return response.Resp().Json(gin.H{"key": "", "expire": 96, "admin": 0}, http.StatusInternalServerError)
			}
		}
		tx.Commit()
	}

	// return response.Resp().Json(gin.H{"token": result.jwtKey, "expire": result.keyExpire, "admin": result.admin}, http.StatusOK)
	return response.Resp().Json(gin.H{"key": privateKey, "expire": keyExpire, "admin": admin, "id": id, "name": account}, http.StatusOK)
}
