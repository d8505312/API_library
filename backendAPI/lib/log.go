package lib

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

//var logger *log.Entry
var Log = logrus.New()

func Initlog() {
	var logfile string = "default.log"
	if Config.GetString("logFile") != "" {
		logfile = Config.GetString("logFile")
	}
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		//Log.Out = file
		Log.SetOutput(io.MultiWriter(file, os.Stdout))
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}

	// 設定日誌格式為json格式
	Log.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	Log.SetLevel(logrus.TraceLevel)
	Log.Info("Log initialized")

}
