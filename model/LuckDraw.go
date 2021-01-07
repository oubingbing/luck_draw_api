package model

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
)

//抽奖类型，1=红包，2=商品
const (
	LUCK_DRAW_TYPE_RED_PAK = 1
	LUCK_DRAW_TYPE_GOODS   = 2
)

//抽奖发布来源，1=平台，2=用户
const (
	LUCK_DRAW_FROM_PLATFORM = 1
	LUCK_DRAW_FROM_USER	    = 2
)

//是否限制参加人数，0=不限制，1=限制
const (
	LUCK_DRAW_LIMIT_JOIN_N = 0
	LUCK_DRAW_LIMIT_JOIN_Y = 1
)

type LuckDrawParam struct {
	Name string `form:"name" json:"name" binding:"required"`
	LimitJoin int32 	 `form:"limit_join" json:"limit_join" binding:"required"` //是否限制参加人数
	JoinLimitNum float32 `form:"join_limit_num" json:"join_limit_num" binding:"required"` //限制参加人数
	ReceiveLimit float32 `form:"receive_limit" json:"receive_limit" binding:"required"` //每人限领数量
	Describe string      `form:"describe" json:"describe" binding:"required"`
	Attachments string   `form:"attachments" json:"attachments" binding:"required"`
	StartAt string    `form:"start_at" json:"start_at" binding:"required"`//活动开始时间
	EndAt string      `form:"end_at" json:"end_at" binding:"required"`//活动截止时间
	RunAt string      `form:"run_at" json:"run_at" binding:"required"`//开奖时间
	ShareTitle string    `form:"share_title" json:"share_title"` //分享标题
	ShareImage string    `form:"share_image" json:"share_image"` //分享图片
}

type LuckDraw struct {
	gorm.Model
	Name string `gorm:"column:name"`
	Type int8   `gorm:"column:type"` //抽奖类型
	FROM int8   `gorm:"column:from"` //发布抽奖的用户类型
	JoinNum int32 		 `gorm:"column:join_num"`   //已参加人数
	LimitJoin int32 	 `gorm:"column:limit_join"`  //是否限制参加人数
	JoinLimitNum float32 `gorm:"column:join_limit_num"` //限制参加人数
	ReceiveLimit float32 `gorm:"column:receive_limit"` //每人限领数量
	Describe string      `gorm:"column:describe"`
	Attachments string   `gorm:"column:attachments"`
	StartAt time.Time    `gorm:"column:start_at"` //活动开始时间
	EndAt time.Time      `gorm:"column:end_at"` //活动截止时间
	RunAt time.Time      `gorm:"column:run_at"` //开奖时间
	ShareTitle string    `gorm:"column:share_title"` //分享标题
	ShareImage string    `gorm:"column:share_image"` //分享图片
}

func (LuckDraw) TableName() string  {
	return "luck_draw"
}

// Value 实现方法
func (p LuckDraw) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p LuckDraw) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), p)
}

func (luckDraw *LuckDraw)Store() (int64,error) {
	createResult := Connect().Create(&luckDraw)
	return createResult.RowsAffected,createResult.Error
}
