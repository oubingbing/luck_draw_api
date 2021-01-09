package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"luck_draw/enums"
	"luck_draw/util"
)

var getCnfErr error = errors.New("配置读取失败")
var connectErr error = errors.New("系统异常")

func Connect() (*gorm.DB,*enums.ErrorInfo) {
	config ,err := util.GetMysqlConfig()
	errorInfo := &enums.ErrorInfo{}
	if err != nil{
		util.Info(fmt.Sprintf("获取数据失败：%v\n",err.Error()))
		errorInfo.Err = getCnfErr
		errorInfo.Code = enums.DB_CONNECT_ERR
		return nil,errorInfo
	}

	db, err := gorm.Open("mysql", config)
	if err != nil {
		util.Info(fmt.Sprintf("连接数据库错误：%v\n",err.Error()))
		errorInfo.Err = connectErr
		errorInfo.Code = enums.DB_CONNECT_ERR
		return nil,errorInfo
	}

	return db,nil
}

