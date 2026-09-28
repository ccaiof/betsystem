package handlers

import (
	"net/http"

	"github.com/betsystem/internal/application/auth"
	"github.com/gin-gonic/gin"
)

func AuthHandler(router *gin.RouterGroup) {
	health := router.Group("/auth")
	{
		health.GET("/me", me)
	}
}

func me(ctx *gin.Context) {
	identity, ok := auth.IdentityFromContext(ctx.Request.Context())

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthenticated",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"subject":    identity.Subject,
		"providerId": identity.ProviderID,
		"roles":      identity.Roles,
	})
}
