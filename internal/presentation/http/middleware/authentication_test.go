package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/betsystem/internal/application/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type fakeVerifier struct {
	identity auth.Identity
	err      error
}

func (f *fakeVerifier) Verify(ctx context.Context, rawToken string) (auth.Identity, error) {
	return f.identity, f.err
}

func TestAuthentication_MissingAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	verifier := &fakeVerifier{}

	router.Use(Authentication(verifier))

	router.GET("/protected", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	require.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}

func TestAuthentication_EmptyBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	verifier := &fakeVerifier{}

	router.Use(Authentication(verifier))

	router.GET("/protected", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer ",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	require.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}

func TestAuthentication_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	verifier := &fakeVerifier{
		err: errors.New("invalid token"),
	}

	router.Use(Authentication(verifier))

	router.GET("/protected", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer invalid-token",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	require.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}

func TestAuthentication_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedIdentity := auth.Identity{
		Subject:    "service-account-provider-acme",
		ProviderID: "provider-acme",
		Roles: []string{
			"provider",
		},
	}

	verifier := &fakeVerifier{
		identity: expectedIdentity,
	}

	router := gin.New()

	router.Use(Authentication(verifier))

	router.GET("/protected", func(ctx *gin.Context) {
		identity, ok := auth.IdentityFromContext(
			ctx.Request.Context(),
		)

		require.True(t, ok)
		require.Equal(
			t,
			expectedIdentity,
			identity,
		)

		ctx.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer valid-token",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)
}
