package keycloak

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVerifier_VerifyProviderToken(t *testing.T) {
	issuer := os.Getenv("KEYCLOAK_ISSUER")

	if issuer == "" {
		t.Fatal("KEYCLOAK_ISSUER is not set")
	}

	clientID := os.Getenv("KEYCLOAK_TEST_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_TEST_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Fatal("Keycloak test credentials are not set")
	}

	config := Config{
		Issuer:   issuer,
		ClientID: os.Getenv("KEYCLOAK_CLIENT_ID"),
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	token := requestClientCredentialsToken(
		t,
		ctx,
		issuer,
		clientID,
		clientSecret,
	)

	verifier, err := NewVerifier(config, ctx)
	require.NoError(t, err)

	identity, err := verifier.Verify(ctx, token)
	require.NoError(t, err)

	require.NotEmpty(t, identity.Subject)
	require.Equal(t, clientID, identity.ProviderID)
	require.Contains(t, identity.Roles, "provider")
}

func requestClientCredentialsToken(
	t *testing.T,
	ctx context.Context,
	issuer string,
	clientID string,
	clientSecret string,
) string {
	t.Helper()

	tokenURL := strings.TrimRight(issuer, "/") +
		"/protocol/openid-connect/token"

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tokenURL,
		strings.NewReader(form.Encode()),
	)
	require.NoError(t, err)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response struct {
		AccessToken string `json:"access_token"`
	}

	require.NoError(
		t,
		json.NewDecoder(resp.Body).Decode(&response),
	)

	require.NotEmpty(t, response.AccessToken)

	return response.AccessToken
}
