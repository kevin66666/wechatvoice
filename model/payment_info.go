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
	OrderNumber     string
	IsPaied         string
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

func GetPaymentQuery(openid string) ([]WechatVoicePaymentInfo, error) {
	db := dbpool.OpenConn()
	defer dbpool.CloseConn(&db)
	list := make([]WechatVoicePaymentInfo, 0)
	err := db.Where("open_id = ?", openid).Where("is_paied = 1").Where("is_solved = ?", "2").Find(&list).Error
	return list, err
}
