package model

import "github.com/jinzhu/gorm"

//奖品类型，1=红包，2=商品
const (
	PRIZE_TYPE_RED_PAK = 1
	PRIZE_TYPE_GOODS   = 2
)

//奖品来源，1=平台，2=用户
const (
	PRIZE_FROM_PLATFORM = 1
	PRIZE_FROM_USER	    = 2
)

type Prize struct {
	gorm.Model
	Name string `gorm:"column:name"`
	Num float32 `gorm:"column:num"`
	Type int8   `gorm:"column:type"`
	FROM int8   `gorm:"column:from"`
	Describe string     `gorm:"column:describe"`
	Attachments string  `gorm:"column:attachments"`
}

func (Prize) TableName() string  {
	return "prize"
}
