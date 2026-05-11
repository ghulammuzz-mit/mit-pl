package controller

import (
	"fmt"
	"mit/platform/internal/envctl/infisical"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects, folders, and secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := context.Background()
		client, err := infisical.New()
		if err != nil {
			return err
		}

		projects, err := client.ListProjects()
		if err != nil {
			return err
		}

		fmt.Println("=== Projects ===")
		for _, p := range projects {
			fmt.Println("-", p.Name, p.ID)
		}

		selected := projects[0].ID
		env, _ := cmd.Flags().GetString("env")

		folders, err := client.ListFolders(selected, env)
		if err != nil {
			return err
		}

		fmt.Println("\n=== Folders ===")
		for _, f := range folders {
			fmt.Println("-", f.Name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
