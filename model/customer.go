package model

import (
	//"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	xorm "github.com/xormplus/xorm"
	//	core "github.com/xormplus/core"
)

type Customer struct {
	CustomerID    string `xorm:"'customerID'"`
	CustomerName  string `xorm:"'customerName'"`
	CustomerPwd   string `xorm:"'customerPwd'"`
	CustomerPhone string `xorm:"'customerPhone'"`
	SelProvince   string `xorm:"'selProvince'"`
	SelCity       string `xorm:"'selCity'"`
	CreateBy      string `xorm:"'createBy'"`
	CreateDate    string `xorm:"'createDate'"`
	UpdateBy      string `xorm:"'updateBy'"`
	UpdateDate    string `xorm:"'updateDate'"`
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

	sqls := "insert into customer (customerID,customerName) values (?,?)"
	res2222, err22 := engin.Exec(sqls, "xiaolun", "2")
	if err22 != nil {
		fmt.Println(err22.Error())
	}
	defer engin.Close()
	fmt.Println(res2222)
	//engin.Insert(sql1, "1", "1", "1")
	for _, k := range res {
		id := k["customerID"]
		fmt.Println("xxxxxxxx")
		fmt.Println(string(id))
	}

}

func GetUserInfoByID(openId string) Customer {
	engin, err := xorm.NewMySQL("mysql", "root:7de2cd9b31@tcp(localhost:3306)/mylawyerfriend")
	if err != nil {
		fmt.Println(err)
	}
	sql := "select * from customer where customerID =?"
	res, err1 := engin.Query(sql, openId)
	if err1 != nil {
		fmt.Println(err1)
	}

	var customer Customer
	// customer.CustomerID = res
	for _, k := range res {
		customer.CustomerID = string(k["customerID"])
		customer.CustomerName = string(k["customerName"])
		customer.CustomerPwd = string(k["customerPwd"])
		customer.CustomerPhone = string(k["customerPhone"])
		customer.SelProvince = string(k["selfProvince"])
		customer.SelCity = string(k["selfCity"])

	}
	defer engin.Close()
	return customer
}

// func SetUserInfo(customerId,customerName,customerPwd,suctomer)

type Lawyer struct {
	LawerId             string
	lawyerPhone         string
	LawyerName          string
	LawyerCertificateNo string
	GroupPhoto          string
	SinglePhoto         string
	SelfProvince        string
	SecCity             string
	LawFirm             string
	GoodAtBusiness      string
	Description         string
	UserId              string //openId
}

func GetLaywerInfos(openId string) Lawyer {
	engin, err := xorm.NewMySQL("mysql", "root:7de2cd9b31@tcp(localhost:3306)/mylawyerfriend")
	if err != nil {
		fmt.Println(err)
	}
	sql := "select * from customer where userID =?"
	res, err1 := engin.Query(sql, openId)
	if err1 != nil {
		fmt.Println(err1)
	}
	var law Lawyer
	for _, k := range res {
		law.LawerId = string(k["lawyerId"])
		law.UserId = string(k["userID"])
		law.GoodAtBusiness = string(k["gooAtBusiness"])
	}
	return law
}

/**
这里其实对用户抽取数据影响不大 主要是说 如果新的用户进来以后  需要到他那边继续注册一次

然而律师要拿到自己的东西的话就必须拿着自己的OpenId 去查询一次数据  然后同步数据过来
*/
