package dbpool

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	MAX_IDLE_CONN_NUM = 0
	MAX_OPEN_CONN_NUM = 100
)

var db gorm.DB
var old_db gorm.DB

func init() {

	var err error

	//设置log模块名称

	//打开数据库获连接
	env := os.Getenv("RUN_ENV")
	if env == "prod" {
		// log.Println("=====>>>连接整死库")
		db, err = gorm.Open("mysql", "root:7de2cd9b31@tcp(localhost:3306)/mylawyerfriend?charset=utf8&parseTime=True&loc=Local")
		//old_db, err = gorm.Open("mysql", "shangqu:riA6Y4y1fEwcg@tcp(wbollalpha.mysql.rds.aliyuncs.com:3306)/wechat?charset=utf8mb4&parseTime=True&loc=Local")
	} else {
		// log.Println("=====>>>连接整死库")
		db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
		//old_db, err = gorm.Open("mysql", "root:llobw123**8@tcp(123.56.11.116:6033)/wechat?charset=utf8mb4&parseTime=True&loc=Local")
	}
	if err != nil {
		log.Println("when connect to mysql:" + err.Error())
		return
	}
	//设置数据库名称单数
	db.SingularTable(true)
	//old_db.SingularTable(true)

	//设置池子大小
	db.DB().SetMaxIdleConns(int(MAX_IDLE_CONN_NUM))
	db.DB().SetMaxOpenConns(int(MAX_OPEN_CONN_NUM))

	db.LogMode(true)
	old_db.LogMode(true)
}

//获取一个可用的数据库连接
func OpenConn() gorm.DB {
	return db
}

//释放数据库连接
func CloseConn(db *gorm.DB) {
}

func GetOldDBConn() gorm.DB {
	return old_db
}
