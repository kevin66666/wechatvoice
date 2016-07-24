package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type WechatVoiceQuestionSettings struct {
	gorm.Model
	Uuid string //PK
	CategoryId string //
	CateGoryName string
	IsSetted string // 是否已经编辑
	PayAmount string //按照分来计算
	PayAmountInt int64 //按照分来计算的int类型
	LawyerFeePercent string //每笔订单律师分润百分比 以string 记录
	UserRedPacketPercent string //每笔订单用户可获取的红包的比例数  (如果律师分润80% 那么 用户红包为剩余20% 再乘以这个数字)
}

func init() {
	info := new(WechatVoiceQuestionSettings)
	info.GetConn().AutoMigrate(&WechatVoiceQuestionSettings{})
}

func (this *WechatVoiceQuestionSettings) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&WechatVoiceQuestionSettings{})
}

func (this *WechatVoiceQuestionSettings) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}

func GetSettingList(startLine,endLine int64)([]WechatVoiceQuestionSettings,error ){
	conn:=dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)

	list :=make([]WechatVoiceQuestionSettings,0)

	err :=conn.Where("uuid is not null").Offset(startLine-1).Limit(endLine-startLine+1).Find(&list).Error
	return list,err
}