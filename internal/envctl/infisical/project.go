package infisical

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	ID   string
	Name string
}

func (c *Client) ListProjects() ([]Project, error) {

	token := c.sdk.Auth().GetAccessToken()

	req, err := http.NewRequest("GET", "https://us.infisical.com/api/v1/projects", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("infisical returned error HTTP %d", resp.StatusCode)
	}

	var parsed struct {
		Projects []Project `json:"projects"`
	}

	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		return nil, err
	}

	return parsed.Projects, nil
}

func (c *Client) SelectProjectInteractive() (string, error) {
	projects, err := c.ListProjects()
	if err != nil {
		return "", err
	}

	fmt.Println("=== Select Project ===")
	for i, p := range projects {
		fmt.Printf("[%d] %s\n", i+1, p.Name)
	}

	fmt.Print("Choose: ")
	var idx int
	fmt.Scanln(&idx)

	return projects[idx-1].ID, nil
}
