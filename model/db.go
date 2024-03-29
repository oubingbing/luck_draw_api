package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"luck_draw/enums"
	"luck_draw/util"
)

type PageParam struct {
	PageNum 		int 		`form:"page_num" json:"page_num" binding:"required"`
	PageSize 		int  		`form:"page_size" json:"page_size" binding:"required"`
	OrderBY 		string 	 	`form:"order_by" json:"order_by" binding:"required"`
	Sort			string 		`form:"sort" json:"sort" binding:"required"` 				//分享图片
	Type			int 		`form:"type" json:"type"` 									//分享图片
	History			int 		`form:"history" json:"history"` 							//分享图片
}

var getCnfErr error = errors.New("配置读取失败")
var connectErr error = errors.New("系统异常")

func Connect() (*gorm.DB,*enums.ErrorInfo) {
	errorInfo := &enums.ErrorInfo{}
	config ,configErr := util.GetMysqlConfig()
	if configErr != nil{
		util.Info(fmt.Sprintf("获取数据失败：%v\n",configErr.Err.Error()))
		return nil,configErr
	}

	db, err := gorm.Open("mysql", config)
	if err != nil {
		util.Info(fmt.Sprintf("连接数据库错误：%v\n",err.Error()))
		errorInfo.Err = connectErr
		errorInfo.Code = enums.DB_CONNECT_ERR
		return nil,errorInfo
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(200)

	cf,_ := util.GetConfig()
	if cf["ENV"] != "prod" {
		//db.LogMode(true)
	}

	return db,nil
}

/**
 * 通用分页
 */
func Page(db *gorm.DB,table string,page *PageParam) *gorm.DB {
	return db.Table(table).Limit(page.PageSize).Offset((page.PageNum-1)*page.PageSize)
}

