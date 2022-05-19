package controller

// import (
// 	. "backendAPI/lib"
// 	"backendAPI/response"
// 	"fmt"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/xuri/excelize/v2"
// )

// type valid struct {
// 	SN     int    `json:"sn"`
// 	Name   string `json:"name"`
// 	Mobile string `json:"mobile"`
// }

// type invalid struct {
// 	SN     int    `json:"sn"`
// 	Name   string `json:"name"`
// 	Mobile string `json:"mobile"`
// 	Msg    string `json:"msg"`
// }

// //驗證發送清單
// func VerifykSendList(context *gin.Context) *response.Response {

// 	file, _, err := context.Request.FormFile("file")
// 	if err != nil {
// 		return response.Resp().Json(gin.H{"msg": "檔案不存在"}, http.StatusNotFound)
// 	}

// 	defer file.Close()

// 	f, err := excelize.OpenReader(file)
// 	if err != nil {
// 		return response.Resp().Json(gin.H{"msg": "無法開啟檔案"}, http.StatusNotFound)
// 	}

// 	cols, err := f.GetCols("發送清單")
// 	if err != nil {
// 		fmt.Println(err)
// 		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
// 	}
// 	var val valid
// 	var vals []valid
// 	var inval invalid
// 	var invals []invalid
// 	tmpArr := make([]string, 0)

// 	for i, col := range cols[1] {
// 		if i > 0 {
// 			mobile := strings.TrimSpace(col)
// 			name := cols[0][i]
// 			tmpArr = append(tmpArr, mobile)

// 			if IsMobile(mobile) {
// 				if IsDuplicate(tmpArr, mobile) {
// 					inval.SN = i
// 					inval.Name = name
// 					inval.Mobile = mobile
// 					inval.Msg = `號碼重複`
// 					invals = append(invals, inval)
// 				} else {
// 					val.SN = i
// 					val.Name = name
// 					val.Mobile = mobile
// 					vals = append(vals, val)
// 				}
// 			} else {
// 				inval.SN = i
// 				inval.Name = name
// 				inval.Mobile = mobile
// 				inval.Msg = `號碼格式有誤`
// 				invals = append(invals, inval)
// 			}
// 		}
// 	}

// 	return response.Resp().Json(gin.H{"valid": vals, "invalid": invals}, http.StatusOK)
// }
