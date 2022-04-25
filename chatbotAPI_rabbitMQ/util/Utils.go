package util

import (
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

//隨機字串
func RandomString(n int) string {
	charSet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(randomString)
}

//檢查空字串
func IsEmpty(str string) bool {
	if str == "" {
		return true
	}
	return false
}

//取得隨機檔案ID
func GetFileID() string {
	// uniqueId := uuid.New()
	// filename := strings.Replace(uniqueId.String(), "-", "", -1)
	filename := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+10000)
	return filename
}

//取得副檔名
func GetFileExt(fileName string) string {
	return strings.ToLower(path.Ext(fileName))
}

//變更檔案名稱
func ChangeFileName(fileName string, newName string) string {
	// fileExt := strings.Split(fileName, ".")[1]
	// filename := fmt.Sprintf("%s.%s", newName, fileExt)
	fileExt := strings.ToLower(path.Ext(fileName))
	filename := fmt.Sprintf("%s%s", newName, fileExt)
	return filename
}

func ConvertGBKToUTF8(str string) string {
	return mahonia.NewDecoder("GBK").ConvertString(str)
}

func ConvertUTF8ToGBK(str string) string {
	return mahonia.NewEncoder("GBK").ConvertString(str)
}
