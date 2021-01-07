package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
)

func CreateLuckDraw(ctx *gin.Context)  {
	var param model.LuckDrawParam
	if err := ctx.ShouldBind(&param); err != nil {
		fmt.Println(err)
		util.ResponseJson(ctx,500,err.Error(),nil)
		return
	}

	service.SaveLuckDraw(&param)

}
