package lib

import (
	"bytes"
	"fmt"
	"math/rand"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//隨機code字串
func GetShortCode(n int) string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, n)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range randomString {
		randomString[i] = charSet[rnd.Intn(len(charSet))]
	}
	return string(randomString)
}

//隨機token字串
func GetTokenString(n int) string {
	charSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, n)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range randomString {
		randomString[i] = charSet[rnd.Intn(len(charSet))]
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
	filename := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+10000)
	return filename
}

//取得副檔名
func GetFileExt(fileName string) string {
	return strings.ToLower(path.Ext(fileName))
}

//變更檔案名稱
func ChangeFileName(fileName string, newName string) string {
	fileExt := strings.ToLower(path.Ext(fileName))
	filename := fmt.Sprintf("%s%s", newName, fileExt)
	return filename
}

func ConvertMapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

func IsLetter(str string) bool {
	var count int
	for _, v := range str {
		if unicode.IsLetter(v) {
			count++
			break
		}
	}
	return count > 0
}

func IsNumber(str string) bool {
	var count int
	for _, v := range str {
		if unicode.IsNumber(v) {
			count++
			break
		}
	}
	return count > 0
}

//時間戳轉時間
func UnixToStr(timeUnix int64, layout string) string {
	timeStr := time.Unix(timeUnix, 0).Format(layout)
	return timeStr
}

//時間轉時間戳
func TimeToUnix(layout, timestr string, duration ...string) int64 {
	t, err := time.Parse(layout, timestr)
	if err != nil {
		Log.Error("[TimeToUnix] ", err)
	}

	var unix int64
	if len(duration) > 0 {
		d := duration[0]
		deadline, _ := time.ParseDuration(d)
		unix = t.Add(deadline).UTC().UnixNano()
	} else {
		unix = t.UTC().UnixNano()
	}

	return unix
}

func StrToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		Log.Error("[StrToInt] ", err)
	}
	return i
}

func StrToFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		Log.Error("[StrToFloat] ", err)
	}
	return f
}

func IsMobile(mobileNum string) bool {
	regular := "^(0|\\+?886)?9\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func IsDuplicate(strArr []string, str string) bool {
	tmpMap := make(map[string]interface{})

	for _, s := range strArr {
		fmt.Println(tmpMap[s])
		if _, ok := tmpMap[str]; ok {
			return true
		}
		tmpMap[s] = s
	}
	return false
}
