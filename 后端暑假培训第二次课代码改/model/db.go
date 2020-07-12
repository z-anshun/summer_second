package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB

func InitDB()  {
	var err error
	db, err := gorm.Open("mysql","mysql","root:123@tcp(127.0.0.1:3306)/summer_course?parseTime=true&charset=utf8&loc=Local")
	if err != nil {
		log.Panicf("Panic while connecting the gorm. Error: %s",err)
	}

	DB = db
	//如果没有对应的表 就创建表
	if !DB.HasTable(&Order{}) {
		DB.CreateTable(&Order{})
	}

	if !DB.HasTable(&Goods{}) {
		DB.CreateTable(&Goods{})
	}

}

 

