package handlers

import "github.com/gin-gonic/gin"

func HealthHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}
