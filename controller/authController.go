package controller

import (
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
)

/**
 * 微信小程序授权登录
 */
func Login(ctx *gin.Context)  {
	loginType,ok := ctx.GetQuery("type")
	if !ok {
		util.ResponseJson(ctx,enums.AUTH_LOGIN_TYPE_ERR,enums.LoginTypeErr.Error(),nil)
		return
	}

	if loginType != "wechat" {

	}

	var loginData enums.WxMiniLoginData
	if err := ctx.ShouldBind(&loginData); err != nil {
		util.ResponseJson(ctx,enums.AUTH_PARAMS_ERROR,err.Error(),nil)
		return
	}

	util.Info("test")

	userInfo,errInfo := service.GetSessionInfo(&loginData)
	if errInfo != nil {
		util.ResponseJson(ctx,errInfo.Code,errInfo.Err.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	token,err := service.UserLogin(db,userInfo)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",token)
	return
}
