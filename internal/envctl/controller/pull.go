package controller

import (
	"fmt"
	fx "mit/platform/internal/envctl/file"
	"mit/platform/internal/envctl/infisical"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull env from Infisical and write to a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := context.Background()
		client, err := infisical.New()
		if err != nil {
			return err
		}

		env, _ := cmd.Flags().GetString("env")
		file, _ := cmd.Flags().GetString("file")

		projectID, err := client.SelectProjectInteractive()
		if err != nil {
			return err
		}

		secretPath, err := client.SelectFolderInteractive(projectID, env)
		if err != nil {
			return err
		}

		secrets, err := client.ListSecrets(projectID, env, secretPath)
		if err != nil {
			return err
		}

		err = fx.Write(file, secrets)
		if err != nil {
			return err
		}

		fmt.Println("Pulled env to", file)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
