package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type Customer struct {
	CustomerId    string
	CustomerName  string
	CustomerPwd   string
	CustomerPhone string
	SelProvince   string
	SelCity       string
	CreateBy      string
	CreateDate    string
	UpdateBy      string
	UpdateDate    string
}

//o-u0Nv8ydozIYnNVzca_C0frKwgI
func (this *Customer) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&Customer{})
}

func (this *Customer) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}
