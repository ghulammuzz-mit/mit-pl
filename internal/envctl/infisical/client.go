package infisical

import (
	"context"
	"fmt"
	"mit/platform/internal/envctl/config"
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
	hostURL := os.Getenv("INFISICAL_HOST_URL")

	// Fallback to ~/.envctl/config [default]
	if clientID == "" || clientSecret == "" {
		prof, err := config.Load("default")
		if err == nil && prof != nil {
			if clientID == "" {
				clientID = prof.ClientID
			}
			if clientSecret == "" {
				clientSecret = prof.ClientSecret
			}
			if hostURL == "" {
				hostURL = prof.HostURL
			}
		}
	}

	if clientID == "" {
		return nil, fmt.Errorf("Infisical client ID not set — run: envctl configure")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("Infisical client secret not set — run: envctl configure")
	}
	if hostURL == "" {
		hostURL = config.DefaultHostURL()
	}

	sdk := infisical.NewInfisicalClient(ctx, infisical.Config{
		SiteUrl:          hostURL,
		AutoTokenRefresh: true,
	})

	_, err := sdk.Auth().UniversalAuthLogin(clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	return &Client{sdk: sdk}, nil
}
