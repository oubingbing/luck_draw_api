package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	JOIN_LOG_STATUS_QUEUE			= 1		//排队中
	JOIN_LOG_STATUS_SUCCESS			= 2		//加入成功
	JOIN_LOG_STATUS_FAIL			= 3		//加入失败
)

type JoinLog struct {
	gorm.Model
	ActivityId 		int64		`gorm:"column:activity_id"` 	//参与活动
	UserId			int64		`gorm:"column:user_id"` 		//参与用户
	Status			int8		`gorm:"column:status"` 			//状态，1=排队中，2=加入成功，3=加入失败
	Remark  		string		`gorm:"column:remark"` 			//备注信息
	JoinedAt 		*time.Time  `gorm:"column:joined_at"` 		//加入的时间
}

func (JoinLog) TableName() string  {
	return "activity_join_log"
}

func (joinLog *JoinLog)Store(db *gorm.DB) (int64,error) {
	result := db.Create(joinLog)
	return result.RowsAffected,result.Error
}

func (joinLog *JoinLog) FindByUserActivity(db *gorm.DB,activityId int64,userId int64) error {
	err := db.Where("activity_id = ?",activityId).Where("user_id = ?",userId).First(joinLog).Error
	fmt.Printf("记录:%v\n,error:%v",db.RecordNotFound(),err)
	return err
}
