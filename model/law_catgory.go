package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type LawCatgory struct {
	gorm.Model
	Uuid        string `sql:"size:32;not null"` //主键
	CategoryName string
}

func init() {
	info := new(LawCatgory)
	info.GetConn().AutoMigrate(&LawCatgory{})
}

func (this *LawCatgory) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&LawCatgory{})
}

func (this *LawCatgory) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}

func GetCateList()([]LawCatgory,error){
	conn :=dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)

	list :=make([]LawCatgory,0)

	err :=conn.Where("uuid is not ?","1").Find(&list).Error

	return list,err
}