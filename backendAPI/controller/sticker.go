package controller

import (
	. "backendAPI/lib"
	"backendAPI/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type stickerPack struct {
	StickerPackID int      `json:"stickerPackID"`
	Count         int      `json:"-"`
	StickerList   []string `json:"stickerList"`
	Ext           string   `json:"ext"`
}

//取得貼圖清單
func Sticker(c *gin.Context) *response.Response {
	var sticker stickerPack
	var stickerList []stickerPack

	qs := "SELECT packID, count, packType FROM sticker WHERE status =?"
	rows, err := DB.Raw(qs, 1).Rows()
	if err != nil {
		Log.Errorf("[Sticker] err:%s", err)
		return response.Resp().Json(gin.H{"msg": "系統異常"}, http.StatusInternalServerError)
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&sticker.StickerPackID,
			&sticker.Count,
			&sticker.Ext,
		)
		sticker.StickerList = getStickerList(sticker.Count)
		stickerList = append(stickerList, sticker)
	}

	if stickerList == nil {
		return response.Resp().Json(gin.H{"msg": "找不到任何sticker"}, http.StatusNotFound)
	}

	return response.Resp().Json(gin.H{"stickerList": stickerList, "prefix": Config.GetString("stickerUrl")}, http.StatusOK)
}

//取得貼圖列表
func getStickerList(count int) []string {
	list := make([]string, count)
	for i := range list {
		list[i] = fmt.Sprint(i + 1)
	}
	list = append(list, "tab")
	list = append(list, "main")
	return list
}
