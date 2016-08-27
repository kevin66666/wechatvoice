package model

import (
	//"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	xorm "github.com/xormplus/xorm"
	//	core "github.com/xormplus/core"
)

type Customer struct {
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
	// engin.SetMapper(core.SameMapper{})
	// custo := new(customer)
	// custo.customerID = "xxxxx"
	// iddd, ddd := engin.Insert(&custo)
	// if ddd != nil {
	// 	fmt.Println(ddd.Error())
	// }
	// fmt.Println(iddd)
	// fmt.Println("xxxxxxxx")

	// fmt.Println("xxxxxxxx")
	// fmt.Println("xxxxxxxx")
	// fmt.Println("xxxxxxxx")
	/*

	 INSERT INTO `wechat_voice_questions` (`category_id_int`,`ask_time`,`answer_id`,`pv`,`is_ranked`,`payment_info`,`important`,`updated_at`,`answer_open_id`,`answer_name`,`asker_head_img`,`payment_info_int`,`appen_question_time`,`have_append_child`,`deleted_at`,`answer_head_img`,`back_rank_point`,`created_at`,`create_time`,`customer_id`,`description`,`uuid`,`voice_path`,`customer_open_id`,`rank_info`,`is_send_to_back`,`solved_time`,`answerd_count`,`is_solved`,`parent_question_id`,`category`,`customer_name`,`is_answerd`,`answerd_time`,`order_number`,`question_type`,`is_paied`,`category_id`) VALUES ('0','2016-08-27 11:03:58','','0','','200','','2016-08-27T11:03:58+08:00','','','','200','0','','<nil>','','0','2016-08-27T11:03:58+08:00','2016-08-27 11:03:58','d06398fc6a6611e600163e105789ad4f','哈哈','e5f8356f6c0211e600163e1057899cd2','','o-u0Nv5Rjxrw2EdmYXqzLXi_uTVo','','','','0','0','','','','0','','1472267038','','','b52725f1670e11e600163e105789bb97')

		**/
	// sql = "update `userinfo` set username=? where id=?"
	sqls := "insert into customer (customerID,customerName) values (?,?)"
	res2222, err22 := engin.Exec(sqls, "xiaolun", "2")
	if err22 != nil {
		fmt.Println(err22.Error())
	}
	fmt.Println(res2222)
	//engin.Insert(sql1, "1", "1", "1")
	for _, k := range res {
		id := k["customerID"]
		fmt.Println("xxxxxxxx")
		fmt.Println(string(id))
	}

}
