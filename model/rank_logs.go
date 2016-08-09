package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type RankInfoLogs struct {
	gorm.Model
	Uuid string
	QuestionId string
	LawyerId string
	LawyerName string
	AskerId string
	AskerName string
	RankInfo string
	RankTime string
	RankPerson string //0 前台 1 后台
}

func init() {
	info := new(RankInfoLogs)
	info.GetConn().AutoMigrate(&RankInfoLogs{})
}

func (this *RankInfoLogs) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&RankInfoLogs{})
}

func (this *RankInfoLogs) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}