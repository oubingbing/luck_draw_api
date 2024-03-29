package model

import (
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"time"
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

//活动分页
type ActivityPageFormat struct {
	ID        		uint
	Name 			string
	GiftId 			int64
	Type 			int8   		 	//活动类型
	FromType 		int32   		 //发布活动的用户类型
	JoinNum 		int32 		   	//已参加人数
	JoinLimitNum 	float32 	 	//限制参加人数
	Attachments 	string
	AttachmentsSli 	[]string
	Status 			int8		 	//活动状态
	Gift			*Gift
	New				int
	CreatedAt		*time.Time
	IsTop			int8
	Number 			string
	DrawType 		int8
}

func (Gift) TableName() string  {
	return "gift"
}

func (gift *Gift)Store(db *gorm.DB) (int64,error) {
	result := db.Create(gift)
	return result.RowsAffected,result.Error
}

func (gift *Gift)First(db *gorm.DB,id int64) (*enums.GiftDetail,error) {
	detail := &enums.GiftDetail{}
	err := db.Table(gift.TableName()).
		Where("deleted_at is null").
		Select("name,user_id,num,type,des,attachments").
		Where("id = ?",id).
		First(detail).Error
	return detail,err
}
