package lib

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func InitConfigure() {
	v := viper.New()
	v.SetConfigName("config") // 設定檔名稱（無後綴）
	v.SetConfigType("yml")    // 設定字尾名 {"1.6以後的版本可以不設定該字尾"}
	v.AddConfigPath("./")     // 設定檔案所在路徑
	v.Set("verbose", true)    // 設定預設引數

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Log.Panicln(" Config file not found; ignore error if desired")
		} else {
			Log.Panicln("Config file was found but another error was produced")
		}
	}

	gin.SetMode(v.GetString("debugMode"))

	// 監控配置和重新獲取配置
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		// if err := viper.Unmarshal(&v); err != nil {
		//     mylog.Println(err.Error())
		// } else {
		//     mylog.Println("config auto reload!")
		// }
	})
	Config = v
}
