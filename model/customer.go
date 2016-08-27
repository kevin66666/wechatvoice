package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	db, err := sql.Open("mysql", "root:7de2cd9b31@tcp(localhost:3306)/mylawyerfriend")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	stmp, _ := db.Prepare("select * from customer where customerID =?")
	row, ss := stmp.Exec("o-u0Nv8ydozIYnNVzca_C0frKwgI")
	bbbb := stmp.QueryRow()
	fmt.Println("ashdkahsdkjahsdjkahsdkjahksdhj")
	fmt.Println(&row)
	fmt.Println("====================================")
	fmt.Println(bbbb)
	if ss != nil {
		fmt.Println(ss.Error())
	}
}
