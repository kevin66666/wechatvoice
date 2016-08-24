package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type WechatVoicePaymentInfo struct {
	gorm.Model
	Uuid            string
	SwiftNumber     string
	QuestionId      string
	MemberId        string
	OpenId          string
	RedPacketAmount string
	LawyerAmount    string
	Left            string
	OrderId         string
}

func init() {
	info := new(WechatVoicePaymentInfo)
	info.GetConn().AutoMigrate(&WechatVoicePaymentInfo{})
}

func (this *WechatVoicePaymentInfo) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&WechatVoicePaymentInfo{})
}

func (this *WechatVoicePaymentInfo) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}
