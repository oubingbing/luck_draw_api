package model

import (
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"time"
)

const (
	JOIN_LOG_STATUS_QUEUE			= 1		//排队中
	JOIN_LOG_STATUS_SUCCESS			= 2		//加入成功
	JOIN_LOG_STATUS_FAIL			= 3		//加入失败
	JOIN_LOG_STATUS_WIN				= 4		//已中奖
	JOIN_LOG_STATUS_LOSE			= 5		//未中奖
)

type JoinLog struct {
	gorm.Model
	ActivityId 		int64		`gorm:"column:activity_id"` 	//参与活动
	UserId			int64		`gorm:"column:user_id"` 		//参与用户
	Status			int8		`gorm:"column:status"` 			//状态，1=排队中，2=加入成功，3=加入失败
	Remark  		string		`gorm:"column:remark"` 			//备注信息
	JoinedAt 		*time.Time  `gorm:"column:joined_at"` 		//加入的时间
}

type JoinLogPage []enums.JoinLogTrans
type JoinLogMemberPage []enums.JoinLogMember

func (JoinLog) TableName() string  {
	return "activity_join_log"
}

func (joinLog *JoinLog)Store(db *gorm.DB) (int64,error) {
	result := db.Create(joinLog)
	return result.RowsAffected,result.Error
}

func (joinLog *JoinLog) FindByUserActivity(db *gorm.DB,activityId int64,userId int64) error {
	err := db.Where("activity_id = ?",activityId).
		Where("user_id = ?",userId).
		Where("status != ?",JOIN_LOG_STATUS_FAIL).
		First(joinLog).Error
	return err
}

func (joinLog *JoinLog) FindById(db *gorm.DB,id string) error {
	err := db.Table(joinLog.TableName()).Where("id = ?",id).First(joinLog).Error
	return err
}

func (joinLog *JoinLog)Update(db *gorm.DB,id uint,data map[string]interface{}) error {
	err := db.Table(joinLog.TableName()).Where("id = ?",id).Updates(data).Error
	return err
}

func (joinLog *JoinLog)LockById(db *gorm.DB,id string) error {
	err := db.Table(joinLog.TableName()).
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?",id).
		First(joinLog).Error

	return err
}

func (joinLog *JoinLog)GetByUserId(db *gorm.DB,userId interface{},status string) (JoinLogPage,error) {
	var page JoinLogPage
	builder := db.Table(joinLog.TableName()).
		Joins("left join activity on activity.id = activity_join_log.activity_id").
		Select("activity_join_log.id,activity_id,user_id,activity_join_log.status,remark,joined_at,activity_join_log.created_at,activity.name,activity.attachments,activity.join_num,activity.join_limit_num,activity.status as activity_status").
		Where("user_id = ?",userId)

	var err error
	if status == "1" {
		err = builder.Where("activity_join_log.status in (?)",[]int8{JOIN_LOG_STATUS_SUCCESS,JOIN_LOG_STATUS_WIN,JOIN_LOG_STATUS_LOSE}).Find(&page).Error
	}else{
		err = builder.Where("activity_join_log.status = ?",status).Find(&page).Error
	}

	return page,err
}

func (joinLog *JoinLog) FindMember(db *gorm.DB,activityId interface{}) (JoinLogMemberPage,error) {
	var page JoinLogMemberPage
	err := db.Table(joinLog.TableName()).
		Joins("left join wechat_user on wechat_user.id = activity_join_log.user_id").
		Select("activity_join_log.id,activity_id,user_id,wechat_user.nick_name,wechat_user.avatar_url").
		Where("activity_id = ?",activityId).
		Where("activity_join_log.status != ?",JOIN_LOG_STATUS_FAIL).
		Find(&page).Error

	return page,err
}

func (joinLog *JoinLog) GetJoinLogByActivityId(db *gorm.DB,activityId uint) ([]JoinLog,error) {
	var joinLogSli []JoinLog
	err := db.Where("activity_id = ?",activityId).
		Where("deleted_at is null").
		Where("status = ?",JOIN_LOG_STATUS_SUCCESS).
		Find(&joinLogSli).Error
	return joinLogSli,err
}
