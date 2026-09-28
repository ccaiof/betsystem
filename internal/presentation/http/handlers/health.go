package handlers

import (
	"net/http"

	"github.com/betsystem/internal/infrastructure/postgres"
	"github.com/gin-gonic/gin"
)

func HealthHandler(router *gin.Engine, postgresHealth *postgres.HealthChecker) {
	health := router.Group("/health")
	{
		health.GET("/live", live)
		health.GET("/ready", ready(postgresHealth))
	}
}

func live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func ready(postgresHealth *postgres.HealthChecker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := postgresHealth.Check(ctx.Request.Context()); err != nil {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}

		// TODO: include health check to SQS

		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}
