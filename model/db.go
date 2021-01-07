package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"luck_draw/util"
)

func Connect() *gorm.DB {
	config ,err := util.GetMysqlConfig()
	fmt.Println(config)
	if err != nil{
		util.Info(fmt.Sprintf("获取数据失败：%v\n",err.Error()))
	}

	db, err := gorm.Open("mysql", config)
	if err != nil {
		util.Info(fmt.Sprintf("连接数据库错误：%v\n",err.Error()))
	}

	return db
}

