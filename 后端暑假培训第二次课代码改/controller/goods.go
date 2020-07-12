package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"summerCourse/model"
	"summerCourse/service"
)

func SelectGoods(ctx *gin.Context) {
	goods := service.SelectGoods()
	ctx.JSON(http.StatusOK, gin.H{
		"status": 200,
		"info":   "success",
		"data": struct {
			Goods []service.Goods `json:"goods"`
		}{goods},
	})
}

//添加商品
func AddGoods(ctx *gin.Context) {
	good := model.Goods{}

	if err := ctx.BindJSON(&good); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 204,
			"info":   "add error",
			"data":   "",
		})
		return
	}
	if err := good.AddGoods(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 204,
			"info":   "db add error",
			"data":   "",
		})
		return
	} else {
		if err := service.AddGoods(service.Goods{good.ID, good.Name, good.Price, good.Num}); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 204,
				"info":   err.Error(),
				"data":   "",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 200,
				"info":   "add success",
				"data":   good,
			})
		}
	}
}
