package main

import (
	"chatbotAPI/controller"
	. "chatbotAPI/lib"
)

// . "pluginAPI/lib"

func main() {
	//讀取設定
	InitConfigure()
	Initlog()
	Initdb()
	InitInflux()
	InitMQ()
	InitMQ2()
	// s := `{"janusMsg":{"chatroomID":"204","msgType":1,"sender":1,"msgContent":"+4謝謝","time":1648203397608000000,"type":2,"format":{},"config":{"id":"yzGrBSCkjGuo2hHZv3wVm","isReply":false,"replyObj":"","currentDate":"2022-03-25","isExpire":false,"isPlay":false,"isRead":true,"msgFunctionStatus":false,"msgMoreStatus":false,"recallPopUp":false,"recallStatus":false}}}`
	// var data []byte = []byte(s)
	// controller.TestAIMessageDate(data)
	controller.ReceiveMQ()
}
