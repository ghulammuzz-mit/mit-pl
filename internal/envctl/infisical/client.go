package infisical

import (
	"context"
	"fmt"
	"os"

	infisical "github.com/infisical/go-sdk"
)

type Client struct {
	sdk infisical.InfisicalClientInterface
}

func New() (*Client, error) {
	ctx := context.Background()

	clientID := os.Getenv("INFISICAL_UNIVERSAL_AUTH_CLIENT_ID")
	clientSecret := os.Getenv("INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET")
	if clientID == "" {
		return nil, fmt.Errorf("INFISICAL_UNIVERSAL_AUTH_CLIENT_ID is required")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET is required")
	}

	sdk := infisical.NewInfisicalClient(ctx, infisical.Config{
		SiteUrl:          "https://us.infisical.com",
		AutoTokenRefresh: true,
	})

	_, err := sdk.Auth().UniversalAuthLogin(clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	return &Client{sdk: sdk}, nil
}
