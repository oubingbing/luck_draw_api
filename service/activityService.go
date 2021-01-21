package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/util"
	"time"
)

var startDateErr 			error 		= errors.New("活动开始日期格式错误")
var endDateErr 				error 		= errors.New("活动截止日期格式错误")
var runDateErr 				error 		= errors.New("活动开奖日期格式错误")
var activityDetailNotFound 	error 		= errors.New("活动详情不存在")
var joinLimit 				error 		= errors.New("活动参与人数达到限制啦")
var saveJoinLogFail 		error 		= errors.New("参加活动失败")
var existsJoinLog	 		error 		= errors.New("您已参加该活动，不可重复参加")
var queryJoinLogDbErr	 	error 		= errors.New("查询出错")

func SaveActivity(db *gorm.DB,param *enums.ActivityCreateParam) (int64,*enums.ErrorInfo) {
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

	for index,_ := range activities {
		activities[index].AttachmentsSli,err = AppendDomain(activities[index].Attachments)
		if err != nil {
			return nil,err
		}
		activities[index].Attachments = ""
	}

	return activities,nil
}

func StrToArr(str string) ([]string,*enums.ErrorInfo) {
	var sli []string
	err := json.Unmarshal([]byte(str),&sli)
	if err != nil {
		return nil,&enums.ErrorInfo{enums.DecodeErr,enums.DECODE_ARR_ERR}
	}

	return sli,nil
}

func AppendDomain(str string) ([]string,*enums.ErrorInfo) {
	config,_ := util.GetConfig()
	domain := config["COS_DOMAIN"]
	sli,err := StrToArr(str)
	if err != nil {
		return nil,err
	}

	for index,_ := range sli {
		sli[index] = domain+"/"+sli[index]
	}

	return sli,nil
}

func ActivityDetail(db *gorm.DB,id string) (*enums.ActivityDetailFormat,*enums.ErrorInfo) {
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

	var parseErr *enums.ErrorInfo
	detail.AttachmentsSli,parseErr = AppendDomain(detail.Attachments)
	if parseErr != nil {
		return nil,parseErr
	}
	detail.Attachments = ""

	detail.ShareImageSli,parseErr = AppendDomain(detail.ShareImage)
	if parseErr != nil {
		return nil,parseErr
	}
	detail.ShareImage = ""

	detail.Gift.AttachmentsSli,parseErr = AppendDomain(detail.Gift.Attachments)
	if parseErr != nil {
		return nil,parseErr
	}
	detail.Gift.Attachments = ""

	return detail,nil
}

/**
 * 进入参与活动队列
 */
func ActivityJoin(db *gorm.DB,id string,userId int64) (*enums.ErrorInfo) {
	activity := &model.Activity{}

	tx := db.Begin()

	//悲观锁
	acNotFound,err := activity.LockById(tx,id)
	if err != nil {
		tx.Rollback()
		util.ErrDetail(enums.ACTIVITY_DETAIL_QUERY_ERR,"活动详情查询错误-"+err.Error(),id)
		return &enums.ErrorInfo{err,enums.ACTIVITY_DETAIL_QUERY_ERR}
	}

	if acNotFound {
		tx.Rollback()
		util.ErrDetail(enums.ACTIVITY_DETAIL_NOT_FOUND,"活动详情不存在-",id)
		return &enums.ErrorInfo{activityDetailNotFound,enums.ACTIVITY_DETAIL_NOT_FOUND}
	}

	if float32(activity.JoinNum) >= activity.JoinLimitNum {
		tx.Rollback()
		return &enums.ErrorInfo{joinLimit,enums.ACTIVITY_JOIN_LIMIT}
	}

	//写入参与日志
	_,joinLogErr := SaveJoinLog(tx,int64(activity.ID),userId)
	if joinLogErr != nil {
		tx.Rollback()
		return joinLogErr
	}
	
	//加入Redis队列进行参与活动

	tx.Commit()

	return nil
}

/**
 * 写入参与日志
 */
func SaveJoinLog(db *gorm.DB,activityId int64,userId int64) (*model.JoinLog,*enums.ErrorInfo) {
	joinLog := &model.JoinLog{}

	err := joinLog.FindByUserActivity(db,activityId,userId)
	if err != nil && !gorm.IsRecordNotFoundError(err){
		util.ErrDetail(
			enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,
			fmt.Sprintf("查询是否重复参与活动出错：%v",err.Error()),
			fmt.Sprintf("activity_id:%v，user_id:%v",activityId,userId))
		return nil,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_QUERY_ERR,Err:queryJoinLogDbErr}
	}

	//如果还有错，那一定是record not found
	if gorm.IsRecordNotFoundError(err) {
		joinLog.ActivityId = activityId
		joinLog.UserId = userId
		joinLog.Status = model.JOIN_LOG_STATUS_FAIL
		joinLog.Remark = ""
		joinLog.JoinedAt = nil

		effect,err := joinLog.Store(db)
		if err != nil {
			util.ErrDetail(
				enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,
				fmt.Sprintf("写入参与日志失败：%v",err.Error()),
				fmt.Sprintf("activity_id:%v，user_id:%v",activityId,userId))
			return nil,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,Err:saveJoinLogFail}
		}

		if effect <= 0 {
			util.ErrDetail(
				enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,
				fmt.Sprintf("写入参与日志失败：%v",effect),
				fmt.Sprintf("activity_id:%v，user_id:%v",activityId,userId))
			return nil,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,Err:saveJoinLogFail}
		}

		return joinLog,nil
	}else{
		return nil,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_REPEAT,Err:existsJoinLog}
	}
}