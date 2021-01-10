package middleware

import (
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/util"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if len(token) <= 0 {
			token = ctx.Query("token")
			if len(token) <= 0 {
				util.ResponseJson(ctx,enums.AUTH_TOKEN_NULL,enums.TokenNull.Error(),nil)
				ctx.Abort()
				return
			}
		}else{
			tokenSlice := strings.Split(token," ")
			token = tokenSlice[1]
			if len(token) <= 0 {
				util.ResponseJson(ctx,enums.AUTH_TOKEN_NULL,enums.TokenNull.Error(),nil)
				ctx.Abort()
				return
			}
		}

		data,err := util.ParseToken(token)
		if err != nil {
			if err == enums.TokenExpired {
				util.ResponseJson(ctx,enums.AUTH_TOKEN_EXPIRED,enums.TokenExpired.Error(),nil)
				ctx.Abort()
				return
			}else{
				util.ResponseJson(ctx,enums.AUTH_USER_PARSE_JWT_ERR,enums.TokenNotValid.Error(),nil)
				ctx.Abort()
				return
			}
		}


		ctx.Set("user_id",data["Id"])
		ctx.Set("open_id",data["OpenId"])

		ctx.Next()

	}
}
