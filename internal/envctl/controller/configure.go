package controller

import (
	"bufio"
	"fmt"
	"mit/platform/internal/envctl/config"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure envctl credentials (~/.envctl/config)",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		existing, _ := config.Load(profile)
		if existing == nil {
			existing = &config.Profile{}
		}

		reader := bufio.NewReader(os.Stdin)

		clientID := prompt(reader, "Infisical Client ID", mask(existing.ClientID))
		if clientID == "" {
			clientID = existing.ClientID
		}

		clientSecret := prompt(reader, "Infisical Client Secret", mask(existing.ClientSecret))
		if clientSecret == "" {
			clientSecret = existing.ClientSecret
		}

		hostURL := prompt(reader, "Infisical Host URL", current(existing.HostURL, config.DefaultHostURL()))
		if hostURL == "" {
			hostURL = existing.HostURL
			if hostURL == "" {
				hostURL = config.DefaultHostURL()
			}
		}

		if err := config.Save(profile, config.Profile{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			HostURL:      hostURL,
		}); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("\nSaved to %s [%s]\n", config.Path(), profile)
		return nil
	},
}

func init() {
	configureCmd.Flags().String("profile", "default", "Profile name")
	rootCmd.AddCommand(configureCmd)
}

func prompt(r *bufio.Reader, label, hint string) string {
	fmt.Printf("%s [%s]: ", label, hint)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func mask(s string) string {
	if s == "" {
		return "None"
	}
	if len(s) <= 4 {
		return "****"
	}
	return "****" + s[len(s)-4:]
}

func current(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
