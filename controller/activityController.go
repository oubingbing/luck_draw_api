package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
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
	if !ok {
		util.ResponseJson(ctx,enums.ACTIVITY_DETAIL_PARAM_ERR,"参数不能为空",nil)
		return
	}

	db,connectErr := model.Connect()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	activity,err := service.ActivityDetail(db ,id)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

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
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	err := service.ActivityJoin(db,id.(string),int64(userId))
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"处理中...",nil)
	return
}
