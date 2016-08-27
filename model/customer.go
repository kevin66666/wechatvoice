package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type Customer struct {
	CustomerID    string
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

func GetCustInfo(id string) Customer {
	db := dbpool.OpenConn()

	defer dbpool.CloseConn(&db)

	var cs Customer
	db.Raw("select * from customer where customerID = ?", id).Find(&cs)
	return cs
}
