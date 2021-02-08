package routers

import (
	"github.com/gin-gonic/gin"
	"luck_draw/controller"
	"luck_draw/middleware"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	notAuth := router.Group("/api")
	notAuth.Use(middleware.NoAuth())
	{
		//微信登录
		notAuth.POST("/login",controller.Login)

		//活动分页
		notAuth.GET("/activity/page",controller.GetActivities)

		//活动详情
		notAuth.GET("/activity/detail",controller.GetDetail)

		//活动参与人员
		notAuth.GET("/activity/member",controller.GetActivityMember)

		//首页广告
		notAuth.GET("/ad/home",controller.AdHome)

		//首页开奖历史广告
		notAuth.GET("/ad/history",controller.AdHome)

		//抽奖详情广告
		notAuth.GET("/ad/detail",controller.AdDetail)

		//消息盒子页面广告
		notAuth.GET("/ad/inbox",controller.AdInbox)

		//消息盒子页面广告
		notAuth.GET("/ad/banner",controller.AdBanner)

		//消息盒子页面广告
		notAuth.GET("/ad/videos",controller.AdVideos)

		//活动类型
		notAuth.GET("/activity/category",controller.ActivityType)

		//获取中奖人员
		notAuth.GET("/activity/wins",controller.GetWins)
	}

	auth := router.Group("/api")
	auth.Use(middleware.Auth())
	{
		//用户信息
		auth.GET("/user/info",controller.GetUserInfo)

		//检测用户登录
		auth.GET("/user/check_login",controller.CheckLogin)

		//用户手机号
		auth.POST("/user/get_phone",controller.GetUserPhone)

		//socket - 获取授权token
		auth.GET("/socket/token",controller.GetSocketToken)

		//新建活动
		auth.POST("/activity/create",controller.CreateActivity)

		//活动参与
		auth.POST("/activity/join",controller.Join)

		//活动参与记录
		auth.GET("/activity/join_log",controller.ActivityLog)

		//分享活动
		auth.POST("/activity/share_join",controller.ShareActivity)

		//新建礼品
		auth.POST("/gift/create",controller.CreateGift)

		//新建地址
		auth.POST("/address/create",controller.CreateAddress)

		//更新地址
		auth.PUT("/address/update",controller.UpdateAddress)

		//删除地址
		auth.DELETE("/address/delete",controller.DeleteAddress)

		//获取地址表信息
		//auth.GET("/address/list",controller.GetAddressList)

		//获取地址表分页
		auth.GET("/address/page",controller.GetAddressPage)

		//获取地址详情
		auth.GET("/address/detail",controller.GetAddressDetail)

		//阅读消息盒子
		auth.PUT("/inbox/read",controller.ReadInbox)

		//消息盒子分页
		auth.GET("/inbox/page",controller.GetInboxPage)

		//未读消息盒子
		auth.GET("/inbox/un_read",controller.GetUnReadInbox)

		//未读消息盒子
		auth.POST("/activity/share",controller.ShareActivity)
	}

	return router
}
