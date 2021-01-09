package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"net/http"
)

type Response struct {
	ErrorCode int `json:"code"`
	ErrorMessage string `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *Response) ResponseError(ctx *gin.Context)  {
	Error(fmt.Sprintf("err_code:%v，err_msg：%v",r.ErrorCode,r.ErrorMessage))
	ctx.JSON(http.StatusOK,r)
}

func (r *Response) ResponseSuccess(ctx *gin.Context)  {
	ctx.JSON(http.StatusOK,r)
}

func ResponseJson(ctx *gin.Context,code int,message string,data interface{})  {
	var res Response
	res.ErrorCode = code
	res.ErrorMessage = message
	res.Data = data

	if code != enums.SUCCESS {
		res.ResponseError(ctx)
	}else{
		res.ResponseSuccess(ctx)
	}
}
