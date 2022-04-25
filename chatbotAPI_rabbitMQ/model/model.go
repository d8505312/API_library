package model

import (
	"encoding/json"
)

//訊息
type Message struct {
	ChatroomID string                 `json:"chatroomID" binding:"required"`
	Sender     int                    `json:"sender" binding:"required"`
	Type       int                    `json:"type" binding:"required"`
	MsgType    *int                   `json:"msgType" binding:"required"`
	MsgContent string                 `json:"msgContent" binding:"required"`
	Time       int64                  `json:"Time" binding:"required"`
	Format     map[string]interface{} `json:"format" binding:"required"`
	Config     map[string]interface{} `json:"config" binding:"required"`
}

type Event struct {
	Plugin string
	Data   map[string]interface{}
}

/*圖片數據 {
	Fileid:      檔案id
	ShowName:    檔案顯示名稱
	Description: 說明
	Width:       圖片寬度
	Height:      圖片高度
	FileSize:    檔案大小
} */
type ImageFileData struct {
	Fileid      int
	ShowName    string
	Description string
	Width       int
	Height      int
	FileSize    int
}

/*貼圖數據 {
	SID:                          表情貼圖id
	PID:                          主題包id
	FileName:                     圖片檔案名稱
	Width:                        圖片寬度
	Height:                       圖片高度
	IsAnimation: 				  是否為動態貼圖
	AnimationImageFileName: 	  動態貼圖檔案名稱
	AnimationHorizontalCutAmount: 動態貼圖水平切割數量
	AnimationVerticalCutAmount:   動態貼圖垂直切割數量
	AnimationTotalAmount:         動態貼圖Frame個數
	AnimationDuration:            動態貼圖每個Frame的持續時間（單位為毫秒）
} */
type StickerData struct {
	SID                          string
	PID                          string
	FileName                     string
	Width                        int
	Height                       int
	IsAnimation                  bool
	AnimationImageFileName       string
	AnimationHorizontalCutAmount int
	AnimationVerticalCutAmount   int
	AnimationTotalAmount         int
	AnimationDuration            int
}

/*錄音數據 {
	Fileid:      檔案id
	ShowName:    檔案顯示名稱
	Description: 說明
	FileSize:    檔案大小（單位：Bytes）
	SoundLength: 音檔長度（單位：秒）
} */
type SoundRecordingFileData struct {
	Fileid      int
	ShowName    string
	Description string
	FileSize    int
	SoundLength int
}

/*文件數據 {
	Fileid:        檔案id
	ShowName:      檔案顯示名稱
	ExtensionName: 副檔名(包括”.”，EX:.doc)
	FileSize:      檔案大小（單位：Bytes）
} */
type DetailFileData struct {
	Fileid        int
	ShowName      string
	ExtensionName string
	FileSize      int
}

/*地圖數據 {
	Longitude:    經度
	Latitude:     緯度
	LocateTime:   定位時間（格式為：yyyy-MM-dd HH:mm:ss +0800）
	LocateSource: 定位來源 (0:Google Map ; 1:Baidu ; 2: Apple)
} */
type MapLocationData struct {
	Longitude    json.Number
	Latitude     json.Number
	LocateTime   string
	LocateSource int
}
