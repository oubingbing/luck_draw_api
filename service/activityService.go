package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"luck_draw/model"
	"time"
)

var startDateErr error = errors.New("活动开始日期格式错误")
var endDateErr error = errors.New("活动截止日期格式错误")
var runDateErr error = errors.New("活动开奖日期格式错误")

func SaveActivity(db *gorm.DB,param *model.ActivityCreateParam) (int64,*enums.ErrorInfo) {
	activity := &model.Activity{
		Name:param.Name,
		GiftId:param.GiftId,
		Type:model.ACTIVITY_TYPE_PHONE_BILL,
		FROM:model.ACTIVITY_FROM_USER,
		LimitJoin:param.LimitJoin,
		JoinLimitNum:param.JoinLimitNum,
		ReceiveLimit:param.ReceiveLimit,
		Describe:param.Describe,
		Attachments:param.Attachments,
		ShareTitle:param.ShareTitle,
		ShareImage:param.ShareImage,
		Status:model.ACTIVITY_STATSUS_TO_RELE,
	}

	var parseErr error
	activity.StartAt,parseErr = time.Parse("2006-01-02 15:04:05",param.StartAt)
	if parseErr != nil {
		return 0,&enums.ErrorInfo{startDateErr,enums.ACTIVITY_START_DATE_ERR}
	}

	activity.EndAt,parseErr = time.Parse("2006-01-02 15:04:05",param.EndAt)
	if parseErr != nil {
		return 0,&enums.ErrorInfo{endDateErr,enums.ACTIVITY_END_DATE_ERR}
	}

	activity.RunAt,parseErr = time.Parse("2006-01-02 15:04:05",param.RunAt)
	if parseErr != nil {
		return 0,&enums.ErrorInfo{runDateErr,enums.ACTIVITY_RUN_DATE_ERR}
	}

	_,err := FirstGiftById(db,activity.GiftId)
	if err != nil {
		return 0,err
	}

	effect,saveErr := activity.Store(db)
	return effect,&enums.ErrorInfo{saveErr,enums.ACTIVITY_SAVE_ERR}
}

func ActivityPage(db * gorm.DB,page *model.ActivityPageParam) (model.AcPage,*enums.ErrorInfo) {
	activity := &model.Activity{}
	activities,err := activity.Page(db,page)
	if err != nil {
		return nil,err
	}

	return activities,nil
}