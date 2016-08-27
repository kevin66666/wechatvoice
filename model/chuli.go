package model

import (
	//"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	xorm "github.com/xormplus/xorm"
	//	core "github.com/xormplus/core"
)

func GetUserInfoByOpenId(openid string) ([]map[string][]byte, error) {
	engin, err := xorm.NewMySQL("mysql", "root:7de2cd9b31@tcp(localhost:3306)/mylawyerfriend")
	if err != nil {
		fmt.Println(err)
	}
	sql := "select * from userinfo where openID =?"
	res, err1 := engin.Query(sql, openid)
	if err1 != nil {
		fmt.Println(err1)
	}
	defer engin.Close()
	return res, err1
}

func GetLawerInfoById(userId string) ([]map[string][]byte, error) {
	engin, err := xorm.NewMySQL("mysql", "root:7de2cd9b31@tcp(localhost:3306)/mylawyerfriend")
	if err != nil {
		fmt.Println(err)
	}
	sql := "select * from  lawyer where userID =?"
	res, err1 := engin.Query(sql, userId)
	if err1 != nil {
		fmt.Println(err1)
	}
	defer engin.Close()
	return res, err1
}
