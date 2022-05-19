package routes

import (
	"backendAPI/controller"
	. "backendAPI/lib"
	"backendAPI/response"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func Load(r *gin.Engine) {
	//跨域設置
	r.Use(CORSMiddleware())

	v1 := r.Group("/")
	{
		v1.GET("/v1/sticker", convert(controller.Sticker))
		v1.GET("/v1/file/:fileid", convert(controller.DownLoadFile))
		v1.GET("/chatbots_autore/:eventId/", convert(controller.Autoresponse))
		v1.GET("/Report", convert(controller.Report))
	}

	v2 := r.Group("/v1/token")
	{
		v2.POST("/", convert(controller.Token))
	}

	// secure v1
	sv1 := r.Group("/v1", JWTAuth())
	// 載入自定義的JWTAuth()中介軟體,在整個sv1的路由組中都生效
	//sv1.Use(jwt.JWTAuth())
	{
		sv1.GET("/", convert(controller.Index))

		sv1.GET("/eventlist", convert(controller.Eventlist))

		sv1.GET("/event/:eventId", convert(controller.EventInfo))

		sv1.GET("/cs", convert(controller.CustomerList))

		sv1.GET("/cs/:eventId", convert(controller.GetEventCustomerList))

		// sv1.PUT("/cs/:eventId", convert(controller.PutEventCustomerList))

		sv1.GET("/events/:accountId", convert(controller.GetCustomerEventList))

		sv1.PUT("/events/:accountId", convert(controller.PutCustomerEventList))

		sv1.POST("/Event", convert(controller.AddEvent))

		sv1.PATCH("/chatroom/:eventId", convert(controller.EditEvent))

		sv1.GET("/chatroomlist/:eventId", convert(controller.ChatRoomList))

		sv1.GET("/livechatroom/:eventId", convert(controller.LiveChatRoom))

		sv1.POST("/history", convert(controller.HistoryMsg))

		sv1.POST("/search", convert(controller.SearchMsg))

		sv1.GET("/chatroom/:chatroomId", convert(controller.GetUserInfo))

		sv1.PATCH("/chatroom", convert(controller.EditUserInfo))

		sv1.POST("/file", convert(controller.UpLoadFile))

		sv1.GET("/accounts", convert(controller.AccountList))

		sv1.POST("/credit", convert(controller.GetCredit))

		sv1.POST("/points", convert(controller.Points))

		sv1.POST("/cancel", convert(controller.CancelReservation))

		sv1.POST("/msgs/:eventId", convert(controller.SendBatchSMS))

		sv1.POST("/msg/:chatroomId", convert(controller.SendSMS))

		sv1.POST("/msg", convert(controller.GetSendMsg))

		sv1.POST("/reply", convert(controller.GetReplyMsg))

		sv1.GET("/admin", convert(controller.GetAdminList))

		sv1.PATCH("/admin", convert(controller.EditAdmin))

		// sv1.POST("/verify", convert(controller.VerifykSendList))

		sv1.GET("/oftens/:accountId", convert(controller.GetOftenSMS))

		sv1.POST("/oftens", convert(controller.AddOftenSMS))

		sv1.PATCH("/oftens", convert(controller.EditOftenSMS))

		sv1.DELETE("/oftens", convert(controller.DeleteOftenSMS))

		sv1.PATCH("/message", convert(controller.EditMessage))

		sv1.POST("/chatbots", convert(controller.Addchatbot))

		sv1.GET("/chatbots/:eventId/:autoId", convert(controller.Getdbchatbot))

		sv1.PATCH("/chatbots", convert(controller.Upchatbot))

		sv1.DELETE("/chatbots", convert(controller.Delchatbot))
	}

	// r.StaticFS("/image", http.Dir("./loadfile"))
}

func convert(f func(*gin.Context) *response.Response) gin.HandlerFunc {

	return func(c *gin.Context) {
		resp := f(c)
		data := resp.GetData()

		switch item := data.(type) {
		case string:
			c.String(resp.GetCode(), item)

		case gin.H:
			c.JSON(resp.GetCode(), item)
		}

		encode, _ := json.Marshal(data)
		Log.Infof("Resp:%s", encode)
	}
}
