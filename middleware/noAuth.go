package middleware

import (
	"github.com/gin-gonic/gin"
)

func NoAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data,err := ParseUser(ctx)
		if err == nil {
			ctx.Set("user_id",data["Id"])
			ctx.Set("open_id",data["OpenId"])
		}

		ctx.Next()
	}
}
