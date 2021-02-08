package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
	"time"
)

var createLDFail error = errors.New("数据保存失败")

/**
 * 新增活动
 */
func CreateActivity(ctx *gin.Context)  {
	var param enums.ActivityCreateParam
	errInfo := &enums.ErrorInfo{}
	if errInfo.Err = ctx.ShouldBind(&param); errInfo.Err != nil {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,errInfo.Err.Error(),nil)
		return
	}

	var effect int64
	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	effect,errInfo = service.SaveActivity(db,&param)
	if errInfo.Err != nil {
		util.ResponseJson(ctx,errInfo.Code,errInfo.Err.Error(),nil)
		return
	}

	if effect <= 0 {
		util.ResponseJson(ctx,enums.ACTIVITY_SAVE_ERR,createLDFail.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",effect)
	return
}

func GetActivities(ctx *gin.Context)  {
	var param model.PageParam
	errInfo := &enums.ErrorInfo{}
	if errInfo.Err = ctx.ShouldBind(&param); errInfo.Err != nil {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,errInfo.Err.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	activities,err := service.ActivityPage(db,&param)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",activities)
	return
}

/**
 * 获取详情
 */
func GetDetail(ctx *gin.Context)  {
	id,ok := ctx.GetQuery("id")
	uid,_:= ctx.Get("user_id")
	if !ok {
		util.ResponseJson(ctx,enums.ACTIVITY_DETAIL_PARAM_ERR,"参数不能为空",nil)
		return
	}

	userId,_ := uid.(float64)

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	activity,err := service.ActivityDetail(db ,id,userId)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	activity.Ad = service.GetAd(enums.AD_TYPE_DETAIL)
	util.ResponseJson(ctx,enums.SUCCESS,"",activity)
	return
}

/**
 * 参与活动
 */
func Join(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
	id,ok := util.Input(ctx,"id")
	if !ok {
		util.ResponseJson(ctx,enums.ACTIVITY_JOIN_PARAM_ERR,"参数不能为空",nil)
		return
	}

	userId,cok := uid.(float64)
	if !cok {
		util.Info(fmt.Sprintf("用户user_id:%v",uid))
		util.ResponseJson(ctx,enums.Auth_TRANS_UID_ERR,enums.UserIdTransErr.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	logId,err := service.ActivityJoin(db,id.(string),int64(userId))
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	/*db2,connectErr2 := model.Connect()
	if connectErr2 != nil {
		util.ResponseJson(ctx,connectErr2.Code,connectErr2.Err.Error(),nil)
		return
	}

	msg,finish := queue.AttemptJoin(db2,logId)
	db2.Close()
	if finish != 0 {
		util.ResponseJson(ctx,enums.FAIL,msg,nil)
		return
	}else{
		util.ResponseJson(ctx,enums.SUCCESS,msg,nil)
		return
	}*/

	util.ResponseJson(ctx,enums.SUCCESS,"排队中...",logId)
	return
}

func ActivityLog(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
	status,_:= ctx.GetQuery("status")

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	result,err := service.GetActivityLog(db,uid,status)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"ok",result)
	return
}

func GetActivityMember(ctx *gin.Context)  {
	activityId,ok:= ctx.GetQuery("activity_id")
	if !ok {
		util.ResponseJson(ctx,enums.ACTIVITY_Id_EMPYT,"参数不能为空",nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	page,err := service.GetJoinLogMember(db,activityId)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"ok",page)
	return
}

func ActivityType(ctx *gin.Context)  {
	data := []map[string]interface{}{
		map[string]interface{}{"type":0,"txt":"全部"},
		//map[string]interface{}{"type":1,"txt":"红包"},
		map[string]interface{}{"type":3,"txt":"话费"},
		map[string]interface{}{"type":2,"txt":"礼品"},
		map[string]interface{}{"type":4,"txt":"虚拟奖品"},
	}
	util.ResponseJson(ctx,enums.SUCCESS,"ok",data)
	return
}

func GetWins(ctx *gin.Context)  {
	activityId,ok:= ctx.GetQuery("activity_id")
	if !ok {
		util.ResponseJson(ctx,enums.ACTIVITY_Id_EMPYT,"参数不能为空",nil)
		return
	}

	var param model.PageParam
	errInfo := &enums.ErrorInfo{}
	if errInfo.Err = ctx.ShouldBind(&param); errInfo.Err != nil {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,errInfo.Err.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	result,err := service.WinMember(db,activityId,&param)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"ok",result)
	return
}

func ShareActivity(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	var hadJoin int64
	joinLog := &model.JoinLog{}
	hadJoin,err := joinLog.CountTodayJoinLog(db,uid)

	if err == nil {
		if hadJoin > 6 {
			//已经超过限制
			util.ResponseJson(ctx,enums.ACTIVITY_JOIN_LIMIT_TIME,enums.ActivityJoinLimit.Error(),nil)
			return
		}
	}

	var ctxNew = context.Background()
	redis := util.NewRedis()
	curtTime := time.Now().Format(enums.DATE_DAY_FORMAT)
	key := enums.ACTIVITY_USER_COUNT+"_"+fmt.Sprintf("%v",uid)+"_"+curtTime
	redisResult := redis.Client.Get(ctxNew,key)
	if redisResult.Err() != nil {
	}
	if len(redisResult.Val()) <= 0 {
		redis.Client.SetEX(ctxNew,key,1,time.Hour*24)
	}else{
		redis.Client.Incr(ctx,key)
	}

	util.ResponseJson(ctx,enums.SUCCESS,"ok",nil)
	return
}
