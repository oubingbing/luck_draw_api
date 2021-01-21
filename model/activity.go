package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"time"
)

//活动类型，1=红包，2=商品
const (
	ACTIVITY_TYPE_RED_PAK 		= 1
	ACTIVITY_TYPE_GOODS   		= 2
	ACTIVITY_TYPE_PHONE_BILL   	= 3
)

//活动发布来源，1=平台，2=用户
const (
	ACTIVITY_FROM_PLATFORM 	= 1
	ACTIVITY_FROM_USER	    = 2
)

//是否限制参加人数，0=不限制，1=限制
const (
	ACTIVITY_LIMIT_JOIN_N = 0
	ACTIVITY_LIMIT_JOIN_Y = 1
)

//活动状态，1=待发布，2=进行中，3=已下架，4=已结束
const (
	ACTIVITY_STATSUS_TO_RELE 	= 1
	ACTIVITY_STATSUS_RUNNING 	= 2
	ACTIVITY_STATSUS_DOWN 		= 3
	ACTIVITY_STATSUS_FINISH 	= 4

)

type Activity struct {
	gorm.Model
	Name 			string 		`gorm:"column:name"`
	GiftId 			int64 		`gorm:"column:gift_id"`
	Type 			int8   		`gorm:"column:type"` 			//活动类型
	FromType 		int8   		`gorm:"column:from_type"` 			//发布活动的用户类型
	JoinNum 		int32 		`gorm:"column:join_num"`   		//已参加人数
	LimitJoin 		int32 	 	`gorm:"column:limit_join"`  	//是否限制参加人数
	JoinLimitNum 	float32 	`gorm:"column:join_limit_num"` 	//限制参加人数
	ReceiveLimit 	float32 	`gorm:"column:receive_limit"` 	//每人限领数量
	Des 			string      `gorm:"column:des"`
	Attachments 	string   	`gorm:"column:attachments"`
	StartAt 		time.Time   `gorm:"column:start_at"` 		//活动开始时间
	EndAt 			time.Time   `gorm:"column:end_at"` 			//活动截止时间
	RunAt 			time.Time   `gorm:"column:run_at"` 			//开奖时间
	Status 			int8		`gorm:"column:status"` 			//活动状态
	ShareTitle 		string    	`gorm:"column:share_title"` 	//分享标题
	ShareImage 		string    	`gorm:"column:share_image"` 	//分享图片
}


type AcPage []ActivityPageFormat

var pageErr error = errors.New("查询出错")

func (Activity) TableName() string  {
	return "activity"
}

func (activity *Activity)Store(db *gorm.DB) (int64,error) {
	createResult := db.Create(activity)
	return createResult.RowsAffected,createResult.Error
}

func (activity *Activity)Page(db *gorm.DB,page *PageParam) (AcPage,*enums.ErrorInfo) {
	var activities AcPage
	err :=  Page(db,activity.TableName(),page).
			Where("status in (?)",[]int8{ACTIVITY_STATSUS_RUNNING,ACTIVITY_STATSUS_FINISH}).
			Select("id,name,gift_id,type,from_type,join_num,attachments,join_limit_num,status").
			Order("id desc").
			Find(&activities).Error
	if err != nil {
		fmt.Printf("数据错误：%v\n",err)
		return nil,&enums.ErrorInfo{pageErr,enums.ACTIVITY_PAGE_ERR}
	}

	return activities,nil
}

func (activity *Activity) Detail(db *gorm.DB,id string) (*enums.ActivityDetailFormat,bool,error,) {
	activityDetail := &enums.ActivityDetailFormat{}
	err := db.Table(activity.TableName()).
		Select("id,name,gift_id,type,from_type,join_num,limit_join,join_limit_num,des,attachments,share_title,share_image,created_at").
		Where("id = ?",id).
		First(activityDetail).Error
	return activityDetail,db.RecordNotFound(),err
}

func (activity *Activity)LockById(db *gorm.DB,id interface{}) error {
	err := db.Table(activity.TableName()).
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?",id).
		First(activity).Error

	return err
}

func (activity *Activity)Update(db *gorm.DB,id uint,data map[string]interface{}) error {
	err := db.Table(activity.TableName()).Where("id = ?",id).Updates(data).Error
	return err
}
