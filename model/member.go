package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)
type MemberInfo struct {
	gorm.Model
	Uuid        string `sql:"size:32;not null"` //主键
	MpId  string `sql:"size:32;not null"`	    //商户ID
	Level       string `sql:"not null"`         //等级制   MsMemberVipLevelInfo的外键 0 游客
	RegistTime  string `sql:"not null"`         //注册时间 (2006-01-02 15:04:05)
	Type        int64  `sql:"not null"`         // 0 注册 听取 1 注册 消费 2 完善个人信息
	TotalAmount string `sql:"not null"`         // 总消费
	OrderCount  int64  `sql:"not null"`         //总单数
	PhoneNumber string // 电话号
	NickName    string `sql:"not null"` //昵称
	HeadImgUrl  string //头像连接
	Score       string // 积分 //预留字段
	Name        string
	Balance     string // 用户余额 用于后期提现

	QqNumber    string
	WeiboAccount string
}

func init() {
	info := new(MemberInfo)
	info.GetConn().AutoMigrate(&MemberInfo{})
}

func (this *MemberInfo) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&MemberInfo{})
}

func (this *MemberInfo) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}