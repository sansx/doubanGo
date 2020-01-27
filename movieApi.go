package main

import (
	_ "./docs"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
"github.com/swaggo/gin-swagger"
"github.com/swaggo/gin-swagger/swaggerFiles"
"net/http"
)

var (
	db *sql.DB
)

// @title dbMovieApi
// @version 0.0.1
// @description  测试
// @BasePath /api/db/
func movieApi() {
	r := gin.New()
	db = connSql()
	defer db.Close()
	// 创建路由组
	v1 := r.Group("/api/v1")
	v1.OPTIONS("/*route", func(g *gin.Context) {
		g.String(http.StatusOK, "")
	})

	movie := v1.Group("/movie")

	movie.GET("/list/:name", sayHello)

	movie.GET("/id/:userId", record)



	// 文档界面访问URL
	// http://127.0.0.1:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}

// @获取指定ID记录
// @Description get record by ID
// @Tags 测试
// @Accept  json
// @Produce json
// @Param   some_id     path    int     true        "userId"
// @Success 200 {string} string	"ok"
// @Router /accounts/{some_id} [get]
func record(c *gin.Context) {
	resDate := ""
	rows, err := db.Query(`select STR_TO_DATE(substring( showDate ,1,10),'%Y-%m-%d') from movieInfo order by showDate ;`)
	checkErr(err)
	num := []string{}
	for rows.Next() {
		var cost string
		err := rows.Scan(&cost)
		checkErr(err)
		num = append(num, cost)
		fmt.Printf("get: %v", cost)
	}
	mJson, _ := json.Marshal(num)
	fmt.Printf("\naaaa: %v\n", resDate)
	c.String(http.StatusOK, "%v", string(mJson))
}

// @你好世界
// @Description 你好
// @Tags 测试
// @Accept  json
// @Produce json
// @Param   name     path    string     true        "name"
// @Success 200 {string} string	"name,helloWorld"
// @Router /accounts/{name} [get]
func sayHello(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, name+",helloWorld")
}
