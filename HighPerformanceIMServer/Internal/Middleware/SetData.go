package Middleware

import (
	"github.com/gin-gonic/gin"
)

func SetData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("id", ctx.Query("id"))
		ctx.Set("uid", ctx.Query("uid"))
		ctx.Set("name", ctx.Query("name"))
	}
}
