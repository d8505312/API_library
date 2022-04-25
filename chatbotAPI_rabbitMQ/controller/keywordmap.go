package controller

import (
	. "chatbotAPI/lib"
	// "strconv"
	// "time"
)

func Keyword() map[string]map[string]interface{} {
	// test_time := TimeMachine()
	datas := map[string]map[string]interface{}{
		"305": map[string]interface{}{
			"name":     "Sue",
			"keywords": "+1,+2,+3",
			"dates":    "2022/03/08",
			"datee":    "2022/04/15",
			"week":     "1,2,3,5",
			"times":    "1000", //23:30
			"timee":    "1410", //19:00
			"reply":    "您好，烏拉拉拉Sue",
		},
		"I7HVSZW": map[string]interface{}{
			"name":     "Bob",
			"keywords": "+4,+5,+6",
			"dates":    "2022/04/21",
			"datee":    "2022/04/25",
			"week":     "456",
			"times":    "590",  //10:50
			"timee":    "1410", //23:30
			"reply":    "您好，烏拉拉拉Bob",
		},
		"I7HVSZW1": map[string]interface{}{
			"name":     "Jake",
			"keywords": "+7,+8,+9",
			"dates":    "2022/04/08",
			"datee":    "2022/04/12",
			"week":     "123",
			"times":    "19:00",
			"timee":    "21:00",
			"reply":    "您好，烏拉拉拉Jake",
		},
		"073": map[string]interface{}{
			"name":     "Tracy",
			"keywords": "+99",
			"dates":    "2022/04/08",
			"datee":    "2022/04/12",
			"week":     "1,2,3,4,5,",
			"times":    "19:00",
			"timee":    "21:00",
			"reply":    "您好，烏拉拉拉Tracy",
		},
	}
	return datas

}

// func TimeMachine() string {
// 	var timeint int
// 	timeint = 580
// 	timeTicker := time.NewTicker(1 * time.Second)

// 	<-timeTicker.C
// 	test_inttted := timeint + 10

// 	timeTicker.Reset(1 * time.Second)
// 	test_time := strconv.Itoa(test_inttted)
// 	Log.Infof("[取更新時間值] %v", test_time)

// 	return test_time
// }
func Getdata(id string, types string) (string, bool) {

	users := Keyword()
	user := ""
	user_interface, exists := users[id][types]
	if exists {
		user := user_interface.(string)
		return user, exists
	}
	Log.Infof("[取值] key: %v%v [輸出] value%v", id, types, user)
	return user, exists

}
