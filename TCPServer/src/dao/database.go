package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	//Open返回数据库对象，并非数据库连接。初始化了连接池，但没有建立数据库连接
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if err != nil{
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(200)	//最大打开的连接数
	db.SetMaxIdleConns(200)	//最大闲置的连接数
	err = db.Ping()			//请求一个连接，验证
	if err != nil{
		log.Fatalln(err)
	}else{
		log.Println("init database successfully...", db)
	}
}










