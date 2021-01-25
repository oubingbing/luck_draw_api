package controller

import (
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
)

func ReadInbox(ctx *gin.Context)  {
	//uid,_:= ctx.Get("user_id") //考虑到性能，暂时不加userID过滤
	id,ok := util.Input(ctx,"id")
	if !ok {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,"id不能为空",nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	err := service.ReadInbox(db,id)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"ok",nil)
	return
}

func GetInboxPage(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
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

	list,err := service.GetInboxList(db,uid,&param);
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"ok",list)
	return
}

func GetUnReadInbox(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	result,err := service.CountInboxUnRead(db,uid);
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"ok",result)
	return
}