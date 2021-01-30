package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/util"
	"math/rand"
	"sort"
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

	config,_ := util.GetConfig()
	domain := config["COS_DOMAIN"]
	for index,item := range activities {
		activities[index].AttachmentsSli,err = AppendDomain(domain,activities[index].Attachments)
		if err != nil {
			return nil,err
		}
		activities[index].Attachments = ""
		if float32(item.JoinNum) > item.JoinLimitNum {
			activities[index].JoinNum = int32(item.JoinLimitNum)
		}
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

func AppendDomain(domain,str string) ([]string,*enums.ErrorInfo) {
	sli,err := StrToArr(str)
	if err != nil {
		return nil,err
	}

	for index,_ := range sli {
		sli[index] = domain+"/"+sli[index]
	}

	return sli,nil
}

func ActivityDetail(db *gorm.DB,id string,userId float64) (*enums.ActivityDetailFormat,*enums.ErrorInfo) {
	activity := &model.Activity{}
	detail,acNotFound,err := activity.Detail(db,id)
	if err != nil {
		return nil,&enums.ErrorInfo{err,enums.ACTIVITY_DETAIL_QUERY_ERR}
	}

	if acNotFound {
		return nil,&enums.ErrorInfo{activityDetailNotFound,enums.ACTIVITY_DETAIL_NOT_FOUND}
	}

	gift := &model.Gift{}
	giftDetail,err := gift.First(db,detail.GiftId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,&enums.ErrorInfo{giftNotFound,enums.GIFT_NOT_FOUND}
		}

		return nil,&enums.ErrorInfo{err,enums.GIFT_GET_DETAIL_ERR}
	}

	detail.Gift = giftDetail

	config,_ := util.GetConfig()
	domain := config["COS_DOMAIN"]
	var parseErr *enums.ErrorInfo
	detail.AttachmentsSli,parseErr = AppendDomain(domain,detail.Attachments)
	if parseErr != nil {
		return nil,parseErr
	}
	detail.Attachments = ""

	detail.ShareImageSli,parseErr = AppendDomain(domain,detail.ShareImage)
	if parseErr != nil {
		return nil,parseErr
	}
	detail.ShareImage = ""

	detail.Gift.AttachmentsSli,parseErr = AppendDomain(domain,detail.Gift.Attachments)
	if parseErr != nil {
		return nil,parseErr
	}
	detail.Gift.Attachments = ""

	//用户如果是登录状态再查询抽奖记录
	fmt.Printf("用户id：%v\n",int64(userId))
	if int64(userId) > 0 {
		joinLog := &model.JoinLog{}
		err = joinLog.FindByUserActivity(db,int64(detail.ID),int64(userId))
		if err == nil {
			detail.ActivityLog = make(map[string]interface{})
			detail.ActivityLog["id"] = joinLog.ID
			detail.ActivityLog["status"] = joinLog.Status
			detail.ActivityLog["remark"] = joinLog.Remark
			detail.ActivityLog["joined_at"] = joinLog.JoinedAt
		}else{
			util.Error(err.Error())
		}
	}

	if detail.JoinNum > int32(detail.JoinLimitNum) {
		detail.JoinNum = int32(detail.JoinLimitNum)
	}

	return detail,nil
}

/**
 * 进入参与活动队列
 */
func ActivityJoin(db *gorm.DB,id string,userId int64) (uint,*enums.ErrorInfo) {
	activity := &model.Activity{}
	tx := db.Begin()

	//悲观锁
	err := activity.LockById(tx,id)
	if err != nil {
		tx.Rollback()
		util.ErrDetail(enums.ACTIVITY_DETAIL_QUERY_ERR,"活动详情查询错误-"+err.Error(),id)
		return 0,&enums.ErrorInfo{err,enums.ACTIVITY_DETAIL_QUERY_ERR}
	}

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		util.ErrDetail(enums.ACTIVITY_DETAIL_NOT_FOUND,"活动详情不存在-",id)
		return 0,&enums.ErrorInfo{activityDetailNotFound,enums.ACTIVITY_DETAIL_NOT_FOUND}
	}

	if float32(activity.JoinNum) >= activity.JoinLimitNum {
		tx.Rollback()
		return 0,&enums.ErrorInfo{joinLimit,enums.ACTIVITY_JOIN_LIMIT}
	}

	//Faker join
	if int(activity.Really) == model.ACTIVITY_REALLY_N {
		fakerUserErr := JoinFakerUser(tx,activity,userId)
		if fakerUserErr != nil {
			return 0,fakerUserErr
		}
	}


	//写入参与日志
	joinLog,joinLogErr := SaveJoinLog(tx,int64(activity.ID),userId,model.JOIN_LOG_STATUS_QUEUE,model.FAKER_N)
	if joinLogErr != nil {
		tx.Rollback()
		return 0,joinLogErr
	}

	//加入队列
	/*var ctx = context.Background()
	redis := util.NewRedis()
	intCmd := redis.Client.LPush(ctx,enums.ACTIVITY_QUEUE,joinLog.ID)
	if intCmd.Err() != nil {
		util.ErrDetail(
			enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,
			enums.ActivityPushQueueErr.Error(),
			fmt.Sprintf("activity_id:%v，user_id:%v",activity.ID,userId))
		return 0,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,Err:enums.ActivityPushQueueErr}
	}*/

	tx.Commit()

	return joinLog.ID,nil
}

func JoinFakerUser(tx *gorm.DB,activity *model.Activity,userId int64) *enums.ErrorInfo {
	ctx := context.Background()
	redis := util.NewRedis()
	defer redis.Client.Close()
	cacheKey := fmt.Sprintf("%v:%v",model.FAKER_USER_KEY,activity.ID)
	intCmd := redis.Client.Get(ctx,cacheKey)
	if intCmd.Err() != nil {
		return &enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
	}

	var fakerUser []int
	parserErr := json.Unmarshal([]byte(intCmd.Val()),&fakerUser)
	if parserErr != nil {
		//解析数据失败
		return &enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
	}

	sort.Ints(fakerUser)
	activityNum := activity.JoinNum
	for i := 0; i <= len(fakerUser) - 1 ; i++ {
		if int(activityNum) == fakerUser[i] {
			fUser := &model.User{}
			userIds,err := GetFakerUser(tx)
			if err != nil {
				return &enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
			}

			rand.Seed(time.Now().UnixNano()+userId+(int64(i)))
			fakerUserId := rand.Intn(int(len(userIds)))

			fmt.Printf("假用户id：%v\n",userIds[fakerUserId].ID)

			//加入Faker
			fmt.Printf("假用户id：%v\n",int64(fUser.ID))
			_,joinLogErr := SaveJoinLog(tx,int64(activity.ID),int64(userIds[fakerUserId].ID),model.JOIN_LOG_STATUS_SUCCESS,model.FAKER_Y)
			if joinLogErr != nil {
				tx.Rollback()
				return joinLogErr
			}

			activityData := make(map[string]interface{})
			activityData["join_num"] = activity.JoinNum+1
			err = activity.Update(tx,activity.ID,activityData)
			if err != nil {
				tx.Rollback()
				util.ErrDetail(enums.ACTIVITY_DEAL_QUEUE_UPDATE_A_ERR,enums.ActivityUpdateJoinNumFailErr.Error(),activity.ID)
				return &enums.ErrorInfo{enums.ActivityUpdateJoinNumFailErr,enums.ACTIVITY_DEAL_QUEUE_UPDATE_A_ERR}
			}
		}

		//判断参加用户是佛已经满人
	}

	return nil
}

/**
 * 写入参与日志
 */
func SaveJoinLog(db *gorm.DB,activityId int64,userId int64,status int8,faker int8) (*model.JoinLog,*enums.ErrorInfo) {
	joinLog := &model.JoinLog{}

	err := joinLog.FindByUserActivity(db,activityId,userId)
	if err != nil && !gorm.IsRecordNotFoundError(err){
		util.ErrDetail(
			enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,
			fmt.Sprintf("查询是否重复参与活动出错：%v",err.Error()),
			fmt.Sprintf("activity_id:%v，user_id:%v",activityId,userId))
		return nil,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_QUERY_ERR,Err:queryJoinLogDbErr}
	}

	//record not found
	if gorm.IsRecordNotFoundError(err) {
		joinLog.ActivityId = activityId
		joinLog.UserId = userId
		joinLog.Status = status
		joinLog.Faker = faker
		joinLog.Remark = ""

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

func GetActivityLog(db *gorm.DB,userId interface{},status string) (model.JoinLogPage,*enums.ErrorInfo) {
	joinLog := &model.JoinLog{}
	result,err := joinLog.GetByUserId(db,userId,status)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,nil
		}else{
			return nil,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_LOG_QUERY_ERR,Err:enums.SystemErr}
		}
	}

	config,_ := util.GetConfig()
	domain := config["COS_DOMAIN"]

	var appendErr *enums.ErrorInfo
	for index,item := range result {
		result[index].AttachmentsSli,appendErr = AppendDomain(domain,item.Attachments)
		if appendErr != nil {
			return nil,appendErr
		}
		result[index].Attachments = ""
	}

	return result,nil
}

func GetJoinLogMember(db *gorm.DB,activityId interface{}) (model.JoinLogMemberPage,*enums.ErrorInfo) {
	joinLog := &model.JoinLog{}
	page,err := joinLog.FindMember(db,activityId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,nil
		}else{
			util.ErrDetail(enums.ACTIVITY_JOIN_LOG_QUERY_MEMBER_ERR,"查询活动参与人员出错",err.Error())
			return nil,&enums.ErrorInfo{enums.SystemErr,enums.ACTIVITY_JOIN_LOG_QUERY_MEMBER_ERR}
		}
	}

	return page,nil
}