package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type OrderPaymentInfo struct {
	gorm.Model
	Uuid string //PK
	CategoryId string //
	CategoryName string

	QuestionId string
	QuestionName string
	OrderNumber string

	UserPaiedAmount string //用户付款数量 单位是分
	UserPaiedAmountInt int64

	LawyerFee string //律师分润数  单位是分
	LawyerFeeInt int64 //律师分润数 int记录

	RedPacketAmount string //用户roll走红包数量
	RedPacketAmountInt int64 //用户roll走红包数量 int

	BalanceAmount string //平台留取数量
	BalanceAmountInt int64 //平台留取Int

	WeixinSwiftNumber string //流水号
	PaiedTime string //付款时间
}

func init() {
	info := new(OrderPaymentInfo)
	info.GetConn().AutoMigrate(&OrderPaymentInfo{})
}

func (this *OrderPaymentInfo) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&OrderPaymentInfo{})
}

func (this *OrderPaymentInfo) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}