package controller

import (
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/service"
	"luck_draw/util"
)

func AdHome(ctx *gin.Context)  {
	adCode := service.GetAd(enums.AD_TYPE_HOME)
	util.ResponseJson(ctx,enums.SUCCESS,"",adCode)
	return
}

func AdHistory(ctx *gin.Context)  {
	adCode := service.GetAd(enums.AD_TYPE_HISTORY)
	util.ResponseJson(ctx,enums.SUCCESS,"",adCode)
	return
}

func AdDetail(ctx *gin.Context)  {
	adCode := service.GetAd(enums.AD_TYPE_DETAIL_CP)
	util.ResponseJson(ctx,enums.SUCCESS,"",adCode)
	return
}

func AdInbox(ctx *gin.Context)  {
	adCode := service.GetAd(enums.AD_TYPE_INBOX)
	util.ResponseJson(ctx,enums.SUCCESS,"",adCode)
	return
}