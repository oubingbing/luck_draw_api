package routers

import (
	"github.com/gin-gonic/gin"
	"newbug/controller"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/wx_min_login",controller.WxMiniLogin)
	}

	return router
}
