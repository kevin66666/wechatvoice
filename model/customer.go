package model

import (
	//"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
)

type Customer struct {
	CustomerID    string `xorm:"customerID"`
	CustomerName  string `xorm:"customerName"`
	CustomerPwd   string `xorm:"customerPwd"`
	CustomerPhone string `xorm:"customerPhone"`
	SelProvince   string `xorm:"selProvince"`
	SelCity       string `xorm:"selCity"`
	CreateBy      string `xorm:"createBy"`
	CreateDate    string `xorm:"createDate"`
	UpdateBy      string `xorm:"updateBy"`
	UpdateDate    string `xorm:"updateDate"`
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
	custo := new(Customer)
	custo.CustomerID = "xxxxx"
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
