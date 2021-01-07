package routers

import (
	"github.com/gin-gonic/gin"
	"luck_draw/controller"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/wx_min_login",controller.WxMiniLogin)
		api.POST("/create/luck_draw",controller.CreateLuckDraw)
	}

	return router
}
