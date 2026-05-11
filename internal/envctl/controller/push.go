package controller

import (
	"fmt"
	"mit/platform/internal/envctl/infisical"
	fx "mit/platform/internal/envctl/file"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push env file to Infisical",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := context.Background()

		client, err := infisical.New()
		if err != nil {
			return err
		}

		env, _ := cmd.Flags().GetString("env")
		file, _ := cmd.Flags().GetString("file")
		yes, _ := cmd.Flags().GetBool("yes")

		projectID, err := client.SelectProjectInteractive()
		if err != nil {
			return err
		}

		secretPath, err := client.SelectFolderInteractive(projectID, env)
		if err != nil {
			return err
		}

		envMap, err := fx.Read(file)
		if err != nil {
			return err
		}

		if !yes {
			var confirm string
			fmt.Print("Push secrets to Infisical? (y/N): ")
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		err = client.UploadEnv(projectID, env, secretPath, envMap)
		if err != nil {
			return err
		}

		fmt.Println("Push completed!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
