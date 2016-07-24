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
