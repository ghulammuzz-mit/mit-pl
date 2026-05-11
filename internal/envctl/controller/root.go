package controller

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mitenvctl",
	Short: "mitenvctl - Manage Infisical environment variables",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("env", "dev", "Environment (dev, stg, prod)")
	rootCmd.PersistentFlags().String("file", ".env", "Env file to write or read")
	rootCmd.PersistentFlags().Bool("yes", false, "Skip confirmation prompts")
}
