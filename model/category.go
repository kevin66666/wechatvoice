package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type Category struct {
	gorm.Model
	Uuid        string `sql:"size:32;not null"` //主键
	CategoryName string //分类名称

}

func init() {
	info := new(Category)
	info.GetConn().AutoMigrate(&Category{})
}

func (this *Category) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&Category{})
}

func (this *Category) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}

func GetCateList(startLine,endLine int64)([]Category,int64,error){
	conn:=dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list :=make([]Category,0)
	var count int64
	var err error
	err = conn.Where("uuid is not null").Find(&list).Count(&count).Error
	err = conn.Where("uuid is not null").Offset(startLine-1).Limit(endLine - startLine+1).Find(&list).Error
	return list,count,err
}
