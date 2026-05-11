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

	clientID := firstNonEmpty(os.Getenv("INFISICAL_UNIVERSAL_AUTH_CLIENT_ID"), embeddedClientID)
	clientSecret := firstNonEmpty(os.Getenv("INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET"), embeddedClientSecret)
	siteURL := firstNonEmpty(os.Getenv("INFISICAL_HOST_URL"), embeddedSiteURL, defaultSiteURL)

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf(
			"envctl was built without embedded Infisical credentials.\n" +
				"Install the official release binary, or set INFISICAL_UNIVERSAL_AUTH_CLIENT_ID " +
				"and INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET to override.")
	}

	sdk := infisical.NewInfisicalClient(ctx, infisical.Config{
		SiteUrl:          siteURL,
		AutoTokenRefresh: true,
	})

	if _, err := sdk.Auth().UniversalAuthLogin(clientID, clientSecret); err != nil {
		return nil, fmt.Errorf("infisical auth failed: %w", err)
	}

	return &Client{sdk: sdk}, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
