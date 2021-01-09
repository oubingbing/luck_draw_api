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
var activityDetailNotFound error = errors.New("活动详情不存在")

func SaveActivity(db *gorm.DB,param *model.ActivityCreateParam) (int64,*enums.ErrorInfo) {
	activity := &model.Activity{
		Name:param.Name,
		GiftId:param.GiftId,
		Type:model.ACTIVITY_TYPE_PHONE_BILL,
		FromType:model.ACTIVITY_FROM_USER,
		LimitJoin:param.LimitJoin,
		JoinLimitNum:param.JoinLimitNum,
		ReceiveLimit:param.ReceiveLimit,
		Des:param.Des,
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

func ActivityPage(db *gorm.DB,page *model.PageParam) (model.AcPage,*enums.ErrorInfo) {
	activity := &model.Activity{}
	activities,err := activity.Page(db,page)
	if err != nil {
		return nil,err
	}

	return activities,nil
}

func ActivityDetail(db *gorm.DB,id string) (*model.ActivityDetailFormat,*enums.ErrorInfo) {
	activity := &model.Activity{}
	detail,acNotFound,err := activity.Detail(db,id)
	if err != nil {
		return nil,&enums.ErrorInfo{err,enums.ACTIVITY_DETAIL_QUERY_ERR}
	}

	if acNotFound {
		return nil,&enums.ErrorInfo{activityDetailNotFound,enums.ACTIVITY_DETAIL_NOT_FOUND}
	}

	gift := &model.Gift{}
	giftDetail,notFound,err := gift.First(db,detail.GiftId)
	if err != nil {
		return nil,&enums.ErrorInfo{err,enums.GIFT_GET_DETAIL_ERR}
	}
	if notFound {
		return nil,&enums.ErrorInfo{giftNotFound,enums.GIFT_NOT_FOUND}
	}
	detail.Gift = giftDetail

	return detail,nil
}