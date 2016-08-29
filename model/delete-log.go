package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type UserDeleteLogs struct {
	gorm.Model
	Uuid        string
	UserOpenId  string
	OrderId     string
	OrderNumber string
}

func init() {
	info := new(UserDeleteLogs)
	info.GetConn().AutoMigrate(&UserDeleteLogs{})
}

func (this *UserDeleteLogs) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&UserDeleteLogs{})
}

func (this *UserDeleteLogs) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}

func GetUserDeletedList(openId string) ([]UserDeleteLogs, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]UserDeleteLogs, 0)
	listErr := conn.Where("user_open_id = ?", openId).Find(&list).Error
	return list, listErr
}
