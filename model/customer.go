package model

import (
	//"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	xorm "github.com/xormplus/xorm"

	core "github.com/xormplus/core"
)

type customer struct {
	customerID    string `xorm:"'customerID'"`
	customerName  string `xorm:"'customerName'"`
	customerPwd   string `xorm:"'customerPwd'"`
	customerPhone string `xorm:"'customerPhone'"`
	selProvince   string `xorm:"'selProvince'"`
	selCity       string `xorm:"'selCity'"`
	createBy      string `xorm:"'createBy'"`
	createDate    string `xorm:"'createDate'"`
	updateBy      string `xorm:"'updateBy'"`
	updateDate    string `xorm:"'updateDate'"`
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
	//
	// engine.SetMapper(core.SameMapper{})
	engin.SetMapper(core.SameMapper{})
	custo := new(customer)
	custo.customerID = "xxxxx"
	iddd, ddd := engin.Insert(&custo)
	if ddd != nil {
		fmt.Println(ddd.Error())
	}
	fmt.Println(iddd)
	fmt.Println("xxxxxxxx")

	fmt.Println("xxxxxxxx")
	fmt.Println("xxxxxxxx")
	fmt.Println("xxxxxxxx")
	//engin.Insert(sql1, "1", "1", "1")
	for _, k := range res {
		id := k["customerID"]
		fmt.Println("xxxxxxxx")
		fmt.Println(string(id))
	}

}
