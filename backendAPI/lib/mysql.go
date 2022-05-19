package lib

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBConn DBConn
var DB *gorm.DB

// type dbconnConfig struct {
// 	DBUser        string
// 	DBPass        string
// 	DBHost        string
// 	DBPort        string
// 	DBDatabase    string
// 	DBMaxLifeTime int
// 	DBMaxConn     int
// 	DBIdleConn    int
// }

// Init Init
func Initdb() {
	DB = connDB(
		Config.GetString("mysql.user"),
		Config.GetString("mysql.pass"),
		Config.GetString("mysql.ip"),
		Config.GetString("mysql.port"),
		Config.GetString("mysql.db"),
		Config.GetInt("mysql.maxlifetime"),
		Config.GetInt("mysql.maxconn"),
		Config.GetInt("mysql.idleconn"),
	)
}

func connDB(user, pass, host, port, database string, lifeTime, maxCon, idle int) *gorm.DB {
	addr := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true",
		user,
		pass,
		host,
		port,
		database,
	)
	fmt.Println(addr)

	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		Log.Panicln(err) // 連不到就直接panic裡服務重起再連
	}

	sqlDB, err := db.DB()
	if err != nil {
		Log.Panicln(err)
	}

	sqlDB.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second) // 每條連線的存活時間
	sqlDB.SetMaxOpenConns(maxCon)                                   // 最大連線數
	sqlDB.SetMaxIdleConns(idle)                                     // 最大閒置連線數

	return db
}
