package middleware

import (
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/util"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data,err := ParseUser(ctx)
		if err != nil {
			util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
			ctx.Abort()
			return
		}

		ctx.Set("user_id",data["Id"])
		ctx.Set("open_id",data["OpenId"])
		ctx.Next()
	}
}

func ParseUser(ctx *gin.Context) (map[string]interface{},*enums.ErrorInfo) {
	token := ctx.GetHeader("Authorization")
	if len(token) <= 0 {
		token = ctx.Query("token")
		if len(token) <= 0 {
			return nil,&enums.ErrorInfo{enums.TokenNull,enums.AUTH_TOKEN_NULL}
		}
	}else{
		tokenSlice := strings.Split(token," ")
		if len(tokenSlice) > 1 {
			token = tokenSlice[1]
		}else{
			return nil,&enums.ErrorInfo{enums.TokenNull,enums.AUTH_TOKEN_NULL}
		}
	}

	data,err := util.ParseToken(token)
	if err != nil {
		if err == enums.TokenExpired {
			return nil,&enums.ErrorInfo{enums.TokenExpired,enums.AUTH_TOKEN_EXPIRED}
		}else{
			return nil,&enums.ErrorInfo{enums.TokenNotValid,enums.AUTH_USER_PARSE_JWT_ERR}
		}
	}

	return data,nil
}
