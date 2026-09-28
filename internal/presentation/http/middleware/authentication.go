package middleware

import (
	"net/http"
	"strings"

	"github.com/betsystem/internal/application/auth"
	"github.com/gin-gonic/gin"
)

func Authentication(verifier auth.Verifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing_authorization",
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_authorization",
			})
			return
		}

		identity, err := verifier.Verify(
			ctx.Request.Context(),
			parts[1],
		)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_token",
			})
			return
		}

		ctx.Request = ctx.Request.WithContext(
			auth.WithIdentity(
				ctx.Request.Context(),
				identity,
			),
		)

		ctx.Next()
	}
}
