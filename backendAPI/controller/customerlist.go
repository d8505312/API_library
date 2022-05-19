package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type customer struct {
	AccountID int32  `json:"accountID"`
	Name      string `json:"name"`
	Admin     int8   `json:"admin"`
	Events    string `json:"events"`
}

//所有客服人員列表
func CustomerList(context *gin.Context) *response.Response {
	if context.GetInt("admin") == Viewer {
		return response.Resp().Json(gin.H{"msg": "沒有權限"}, http.StatusUnauthorized)
	}

	var cs customer
	var csList []customer

	qs := `SELECT a.accountID, a.account, a.admin, IFNULL((SELECT GROUP_CONCAT(e.name SEPARATOR '/') FROM accountEvent ae, event e WHERE e.eventID = ae.eventID AND a.accountID = ae.accountID), "") AS events FROM account a WHERE a.customerID = ? AND a.admin < ?`
	rows, err := DB.Raw(qs, context.GetInt("customer"), Manager).Rows()
	if err != nil {
		Log.Error("[CustomerList] ", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&cs.AccountID,
			&cs.Name,
			&cs.Admin,
			&cs.Events,
		)
		csList = append(csList, cs)
	}

	if csList == nil {
		return response.Resp().Json(gin.H{"msg": "沒有客服人員"}, http.StatusNotFound)
	}

	return response.Resp().Json(gin.H{"cslist": csList}, http.StatusOK)
}
