package main

import (
	. "backendAPI/lib"
	"backendAPI/routes"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func main() {
	//讀取設定
	InitConfigure()
	Initlog()
	Initdb()
	InitInflux()

	router := gin.Default()

	//主要程式入口
	routes.Load(router)

	// router.Use(TlsHandler())
	panic(router.RunTLS(":"+Config.GetString("port"), Config.GetString("cert"), Config.GetString("key")))
	// Log.Panicln(router.Run(":" + Config.GetString("port")))
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + Config.GetString("port"),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
