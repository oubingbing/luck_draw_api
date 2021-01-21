package routers

import (
	"github.com/gin-gonic/gin"
	"luck_draw/controller"
	"luck_draw/middleware"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	notAuth := router.Group("/api")
	{
		//微信登录
		notAuth.POST("/login",controller.Login)

		//活动分页
		notAuth.GET("/activity/page",controller.GetActivities)

		//活动详情
		notAuth.GET("/activity/detail",controller.GetDetail)
	}

	auth := router.Group("/api")
	auth.Use(middleware.Auth())
	{
		//用户登录检测
		auth.GET("/user/check_login",controller.CheckLogin)
		//用户信息
		auth.GET("/user/info",controller.GetUserInfo)


		//新建活动
		auth.POST("/activity/create",controller.CreateActivity)

		//活动参与
		auth.POST("/activity/join",controller.Join)

		//新建礼品
		auth.POST("/gift/create",controller.CreateGift)

	}

	return router
}
