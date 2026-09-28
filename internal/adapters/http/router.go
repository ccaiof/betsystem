package http

import (
	"github.com/betsystem/internal/adapters/http/handlers"
	"github.com/betsystem/internal/infrastructure/postgres"
	"github.com/gin-gonic/gin"
)

func NewRouter(postgresHealth *postgres.HealthChecker) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers.HealthHandler(router, postgresHealth)

	return router
}
