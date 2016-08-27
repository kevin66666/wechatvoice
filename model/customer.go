package model

import (
	//"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
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

func GetInfo() {
	engin, err := xorm.NewMySQL("mysql", "root:7de2cd9b31@tcp(localhost:3306)/mylawyerfriend")
	if err != nil {
		fmt.Println(err)
	}
	sql := "select * from customer where customerID =?"
	res, err1 := engin.Query(sql, "o-u0Nv8ydozIYnNVzca_C0frKwgI")
	if err1 != nil {
		fmt.Println(err1)
	}
	// engin.Insert()
	sql1 := "insert into customer where customerID =?,customerName =?,selCity =?"
	engin.Insert(sql1, "1", "1", "1")
	for _, k := range res {
		id := k["customerID"]
		fmt.Println("xxxxxxxx")
		fmt.Println(string(id))
	}

}
