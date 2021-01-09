package routers

import (
	"github.com/gin-gonic/gin"
	"luck_draw/controller"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		//微信登录
		api.POST("/wx_min_login",controller.WxMiniLogin)

		//新建活动
		api.POST("/activity/create",controller.CreateActivity)
		//活动分页
		api.GET("/activity/page",controller.GetActivities)

		//新建礼品
		api.POST("/gift/create",controller.CreateGift)

	}

	return router
}
