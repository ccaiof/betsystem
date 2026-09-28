package http

import (
	"github.com/betsystem/internal/application/auth"
	"github.com/betsystem/internal/infrastructure/postgres"
	"github.com/betsystem/internal/presentation/http/handlers"
	"github.com/betsystem/internal/presentation/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(postgresHealth *postgres.HealthChecker, verifier auth.Verifier) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers.HealthHandler(router, postgresHealth)

	protected := router.Group("/")
	protected.Use(middleware.Authentication(verifier))
	handlers.AuthHandler(protected)

	return router
}
