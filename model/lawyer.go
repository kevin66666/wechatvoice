package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type LawyerInfo struct {
	gorm.Model
	Uuid            string `sql:"size:32;not null"` //主键
	MpId            string `sql:"size:32;not null"` //商户ID
	Level           string `sql:"not null"`         //等级制   MsMemberVipLevelInfo的外键 0 游客
	RegistTime      string `sql:"not null"`         //注册时间 (2006-01-02 15:04:05)
	TotalAmount     string `sql:"not null"`         // 总消费
	OrderCount      int64  `sql:"not null"`         //总单数
	PhoneNumber     string // 电话号
	NickName        string `sql:"not null"` //昵称
	HeadImgUrl      string //头像连接
	Score           string // 积分 //预留字段
	Name            string
	Balance         string // 用户余额 用于后期提现
	OpenId          string
	QqNumber        string
	WeiboAccount    string
	FirstCategory   string //lawcatgoryID FK
	FirstCategoryId string
	SecondCategory  string //lawcatgoryID FK
	// SecondCa
	ThridCategory string //lawcatgoryID FK

	RankFirst  int64
	RankSecond int64
	RankThird  int64
	RankFouth  int64
	RankLast   int64
	Cet        string
	GroupPhoto string
	LawFirm    string
	Desc       string
}

func init() {
	info := new(LawyerInfo)
	info.GetConn().AutoMigrate(&LawyerInfo{})
}

func (this *LawyerInfo) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&LawyerInfo{})
}

func (this *LawyerInfo) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}

func GetLayersFinanceQueryInfo(startLine, endLine int64, name string) ([]LawyerInfo, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]LawyerInfo, 0)
	var count int64
	query := conn.Where("uuid is not null")
	if name != "" {
		query = query.Where("nick_name LIKE ?", "%"+name+"%").Or("name LIKE ?", "%"+name+"%").Or("phone_number LIKE ?", "%"+name+"%")
	}
	err := query.Find(&list).Count(&count).Error

	err = query.Offset(startLine - 1).Limit(endLine - startLine + 1).Find(&list).Error

	return list, count, err
}
