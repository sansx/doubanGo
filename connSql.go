package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

//数据库配置
const (
	userName = "root"
	password = "a88869321"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "learnsql"
)

func connSql() *sql.DB {
	//, "?charset=utf8"
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8mb4"}, "")
	db, err := sql.Open("mysql", path)
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
