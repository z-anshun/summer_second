package main

import (
	"github.com/gin-gonic/gin"
	"summerCourse/controller"
	"summerCourse/model"
	"summerCourse/service"
)

func main() {
	//初始化数据库
	model.InitDB()
	//初始化服务
	service.InitService()
	//router
	r := gin.Default()
	r.GET("/getGoods", controller.SelectGoods)
	r.POST("/order", controller.MakeOrder)
	r.POST("/add",controller.AddGoods)
	r.Run(":8080")

}



