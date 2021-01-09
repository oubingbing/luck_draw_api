package controller

import (
	"github.com/gin-gonic/gin"
	"luck_draw/util"
)

type WxMiniLoginData struct {
	Iv string `form:"iv" json:"iv" binding:"required"`
	Code string `form:"code" json:"code" binding:"required"`
	EncryptedData string `form:"encrypted_data" json:"encrypted_data" binding:"required"`
}

/**
 * 微信小程序授权登录
 */
func WxMiniLogin(ctx *gin.Context)  {
	var loginData WxMiniLoginData
	if err := ctx.ShouldBind(&loginData); err != nil {
		util.ResponseJson(ctx,500,err.Error(),nil)
		return
	}
}
