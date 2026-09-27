package keycloak

import "os"

type Config struct {
	Issuer   string
	ClientID string
}

func NewConfig() Config {
	return Config{
		Issuer:   os.Getenv("KEYCLOAK_ISSUER"),
		ClientID: os.Getenv("KEYCLOAK_CLIENT_ID"),
	}
}
