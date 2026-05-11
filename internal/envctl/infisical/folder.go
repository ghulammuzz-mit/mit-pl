package infisical

import (
	"fmt"

	infisical "github.com/infisical/go-sdk"
	"github.com/infisical/go-sdk/packages/models"
)

func (c *Client) ListFolders(projectID, env string) ([]models.Folder, error) {
	return c.sdk.Folders().List(infisical.ListFoldersOptions{
		ProjectID:   projectID,
		Environment: env,
		Path:        "/",
	})
}

func (c *Client) SelectFolderInteractive(projectID, env string) (string, error) {
	folders, err := c.ListFolders(projectID, env)
	if err != nil {
		return "", err
	}

	fmt.Println("=== Select App (Folder) ===")
	for i, f := range folders {
		fmt.Printf("[%d] %s\n", i+1, f.Name)
	}

	fmt.Print("Choose: ")
	var idx int
	fmt.Scanln(&idx)

	return "/" + folders[idx-1].Name, nil
}
