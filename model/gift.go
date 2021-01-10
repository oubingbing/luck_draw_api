package model

import (
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
)

//奖品类型，1=红包，2=商品，3=话费
const (
	GIFT_TYPE_RED_PAK 		= 1
	GIFT_TYPE_GOODS   		= 2
	GIFT_TYPE_PHONE_BILL 	= 3
)

//奖品来源，1=平台，2=用户
const (
	GIFT_FROM_PLATFORM 		= 1
	GIFT_FROM_USER	    	= 2
)

//奖品状态，1=上架，2=下架，下架不可用
const (
	GIFT_STATUS_UP			= 1
	GIFT_STATUS_DOWN		= 2
)

type Gift struct {
	gorm.Model
	Name 		string 		`gorm:"column:name"`
	UserId 		int 		`gorm:"column:user_id"`
	Num 		float32 	`gorm:"column:num"`
	Type 		int8   		`gorm:"column:type"`			//奖品类型，1=红包，2=商品，3=话费
	FROM        int8   		`gorm:"column:from"`  			//奖品来源，1=平台，2=用户
	STATUS      int8   		`gorm:"column:status"`			//奖品状态，1=上架，2=下架，下架不可用
	Des    		string      `gorm:"column:des"`
	Attachments string  	`gorm:"column:attachments"`
}

func (Gift) TableName() string  {
	return "gift"
}

func (gift *Gift)Store(db *gorm.DB) (int64,error) {
	result := db.Create(gift)
	return result.RowsAffected,result.Error
}

func (gift *Gift)First(db *gorm.DB,id int64) (*enums.GiftDetail,bool,error) {
	detail := &enums.GiftDetail{}
	notFound := db.Table(gift.TableName()).
		Select("name,user_id,num,type,des,attachments").
		Where("id = ?",id).
		First(detail).
		RecordNotFound()
	return detail,notFound,db.Error
}
