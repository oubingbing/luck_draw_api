package controller

import (
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/util"
)

func AdHome(ctx *gin.Context)  {
	adCode := "adunit-cfc0b8b293d531df"
	util.ResponseJson(ctx,enums.SUCCESS,"",adCode)
	return
}