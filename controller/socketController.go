package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/service"
	"luck_draw/util"
)

func GetSocketToken(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
	token,err := service.GetSocketToken(fmt.Sprintf("%v",uid))
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	config,_ := util.GetConfig()
	mp := make(map[string]interface{})
	mp["token"] = token
	mp["domain"] = "wss://"+config["SOCKET_DOMAIN"]

	util.ResponseJson(ctx,enums.SUCCESS,"ok",mp)
	return
}
