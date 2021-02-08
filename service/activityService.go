package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/util"
	"math/rand"
	"sort"
	"strconv"
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
	today := time.Now().Format(enums.DATE_DAY_FORMAT)
	for index,item := range activities {
		activities[index].AttachmentsSli,err = AppendDomain(domain,activities[index].Attachments)
		if err != nil {
			return nil,err
		}
		activities[index].Attachments = ""
		if float32(item.JoinNum) > item.JoinLimitNum {
			activities[index].JoinNum = int32(item.JoinLimitNum)
		}

		fmt.Println(activities[index].CreatedAt.Format(enums.DATE_DAY_FORMAT),today)
		if activities[index].CreatedAt.Format(enums.DATE_DAY_FORMAT) == today {
			activities[index].New = 1
		}else{
			activities[index].New = 0
		}
	}

	//获取置顶
	/*tops := []model.ActivityPageFormat{}
	if int(page.Type) == int(0) && page.PageNum == 1 {
		var ctx = context.Background()
		redis := util.NewRedis()
		defer func() {
			redis.Client.Close()
			db.Close()
		}()
		cmd := redis.Client.Get(ctx,model.TOP_ACTIVITY)
		if cmd.Err() == nil {
			if len(cmd.Val()) > 0 {
				//转化数据
				if json.Unmarshal([]byte(cmd.Val()),&tops) != nil {
					//记录错误
				}
			}
		}else{
			//记录错误
		}
	}
	if len(tops) > 0 {
		return append(tops,activities...),nil
	}*/

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
	if len(str) <= 0 {
		return []string{},nil
	}
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

	if len(detail.ShareImage) > 10 {
		detail.ShareImageSli,parseErr = AppendDomain(domain,detail.ShareImage)
		if parseErr != nil {
			return nil,parseErr
		}
		detail.ShareImage = ""
	}

	detail.Gift.AttachmentsSli,parseErr = AppendDomain(domain,detail.Gift.Attachments)
	if parseErr != nil {
		return nil,parseErr
	}
	detail.Gift.Attachments = ""

	//用户如果是登录状态再查询抽奖记录
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

	var hadJoin int64
	joinLog := &model.JoinLog{}
	hadJoin,err := joinLog.CountTodayJoinLog(db,userId)

	if err == nil {
		if hadJoin >= 5 {
			//已经超过限制
			return 0,&enums.ErrorInfo{enums.ActivityJoinLimit,enums.ACTIVITY_JOIN_LIMIT_TIME}
		}

		//ACTIVITY_USER_COUNT
		//加入队列
		//redisCount := 0

		if hadJoin >= 2 {
			//已经超过限制
			var ctx = context.Background()
			redis := util.NewRedis()
			curtTime := time.Now().Format(enums.DATE_DAY_FORMAT)
			key := enums.ACTIVITY_USER_COUNT+"_"+fmt.Sprintf("%v",userId)+"_"+curtTime
			redisResult := redis.Client.Get(ctx,key)
			if len(redisResult.Val()) <= 0 {
				return 0,&enums.ErrorInfo{enums.ActivityJoinLimitShare,enums.ACTIVITY_JOIN_LIMIT_TIME}
			}else{
				redisCount,_:=strconv.Atoi(redisResult.Val())
				fmt.Println("Redis计数")
				fmt.Println(redisCount)
				if redisCount <= 0 {
					return 0,&enums.ErrorInfo{enums.ActivityJoinLimitShare,enums.ACTIVITY_JOIN_LIMIT_TIME}
				}else{
					redis.Client.Decr(ctx,key)
				}
			}
		}
	}

	//悲观锁
	err = activity.LockById(db,id)
	if err != nil {
		util.ErrDetail(enums.ACTIVITY_DETAIL_QUERY_ERR,"活动详情查询错误-"+err.Error(),id)
		return 0,&enums.ErrorInfo{enums.NetErr,enums.ACTIVITY_DETAIL_QUERY_ERR}
	}

	if err == gorm.ErrRecordNotFound {
		util.ErrDetail(enums.ACTIVITY_DETAIL_NOT_FOUND,"活动详情不存在-",id)
		return 0,&enums.ErrorInfo{activityDetailNotFound,enums.ACTIVITY_DETAIL_NOT_FOUND}
	}

	if float32(activity.JoinNum) >= activity.JoinLimitNum {
		return 0,&enums.ErrorInfo{joinLimit,enums.ACTIVITY_JOIN_LIMIT}
	}

	//写入参与日志
	joinLog,joinLogErr := SaveJoinLog(db,int64(activity.ID),userId,model.JOIN_LOG_STATUS_QUEUE,model.FAKER_N)
	if joinLogErr != nil {
		return 0,joinLogErr
	}

	//加入队列
	var ctx = context.Background()
	redis := util.NewRedis()
	intCmd := redis.Client.LPush(ctx,enums.ACTIVITY_QUEUE,joinLog.ID)
	if intCmd.Err() != nil {
		util.ErrDetail(
			enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,
			enums.ActivityPushQueueErr.Error(),
			fmt.Sprintf("activity_id:%v，user_id:%v",activity.ID,userId))
		return 0,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_SAVE_LOG_FAIL,Err:enums.ActivityPushQueueErr}
	}

	return joinLog.ID,nil
}

func JoinFakerUser(tx *gorm.DB,activity *model.Activity,userId int64) *enums.ErrorInfo {
	ctx := context.Background()
	redis := util.NewRedis()
	defer redis.Client.Close()
	cacheKey := fmt.Sprintf("%v:%v",model.FAKER_USER_KEY,activity.ID)
	intCmd := redis.Client.Get(ctx,cacheKey)
	if intCmd.Err() != nil {
		util.Error(fmt.Sprintf("这个faker 用户数组不存在 %v:%v",model.FAKER_USER_KEY,activity.ID))
		return nil
	}

	var fakerUser []int
	parserErr := json.Unmarshal([]byte(intCmd.Val()),&fakerUser)
	if parserErr != nil {
		//解析数据失败
		return &enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
	}

	activityNum := activity.JoinNum
	sort.Ints(fakerUser)
	activityNewNum := activity.JoinNum
	for i := 0; i <= len(fakerUser) - 1 ; i++ {
		if int(activityNum) == fakerUser[i] {
			userIds,err := GetFakerUser(tx)
			if err != nil {
				fmt.Println(intCmd.Err())
				return &enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
			}

			rand.Seed(time.Now().UnixNano()+userId+(int64(i)))
			fakerUserId := rand.Intn(int(len(userIds)))

			//加入Faker
			_,joinLogErr := SaveJoinLog(tx,int64(activity.ID),int64(userIds[fakerUserId].ID),model.JOIN_LOG_STATUS_SUCCESS,model.FAKER_Y)
			if joinLogErr != nil {
				//tx.Rollback()
				return joinLogErr
			}

			activityNewNum += 1
		}

		//判断参加用户是佛已经满人
	}

	activityData := make(map[string]interface{})
	activityData["join_num"] = activityNewNum
	err := activity.Update(tx,activity.ID,activityData)
	if err != nil {
		//tx.Rollback()
		util.ErrDetail(enums.ACTIVITY_DEAL_QUEUE_UPDATE_A_ERR,enums.ActivityUpdateJoinNumFailErr.Error(),activity.ID)
		return &enums.ErrorInfo{enums.ActivityUpdateJoinNumFailErr,enums.ACTIVITY_DEAL_QUEUE_UPDATE_A_ERR}
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

	ip := util.GetLocalIP()
	fmt.Printf("IP地址：%v\n",ip)
	util.Error(fmt.Sprintf("IP地址：%v\n",ip))

	//record not found
	if gorm.IsRecordNotFoundError(err) {
		joinLog.ActivityId = activityId
		joinLog.UserId = userId
		joinLog.Status = status
		joinLog.Faker = faker
		joinLog.IP = ip
		joinLog.Remark = ""
		joinLog.OrderId = fmt.Sprintf("%v%v",time.Now().UnixNano(),userId)

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

func WinMember(db *gorm.DB,activityId interface{},page *model.PageParam) (model.JoinLogMemberPage,*enums.ErrorInfo) {
	joinLog := &model.JoinLog{}
	list,err := joinLog.Wins(db,activityId,page)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,nil
		}else{
			util.ErrDetail(enums.ACTIVITY_JOIN_LOG_QUERY_ERR,"查询获奖者时数据库发生错误",err.Error())
			return nil,&enums.ErrorInfo{Code:enums.ACTIVITY_JOIN_LOG_QUERY_ERR,Err:enums.SystemErr}
		}
	}

	for j,item := range list {
		newName := ""
		encodeName, isBase4 := base64.StdEncoding.DecodeString(item.NickName)
		if isBase4 == nil {
			item.NickName = string(encodeName)
		}

		if len(item.NickName) > 0 {
			str := []rune(item.NickName)
			if len(str) == 1 {
				newName += "*"
			}else{
				for i := 0; i < len(str); i++ {
					if i == 0 {
						newName += string(str[i])
					}else if len(str) == 2 {
						newName += "*"
					}else if i == (len(str)-1){
						newName += string(str[i])
					}else{
						newName += "*"
					}
				}
			}
		}
		list[j].NickName = newName
	}

	return list,nil
}