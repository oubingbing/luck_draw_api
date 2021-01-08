package controller

import (
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
)

func CreateGift(ctx *gin.Context)  {
	var param model.GiftParam
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

	var effect int64
	userId := 1
	effect,errInfo = service.SaveGift(db,userId,&param)
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
