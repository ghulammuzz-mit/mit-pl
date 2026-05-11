package infisical

import (
	infisical "github.com/infisical/go-sdk"
)

func (c *Client) ListSecrets(projectID, env, path string) ([]infisical.Secret, error) {
	return c.sdk.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:          projectID,
		Environment:        env,
		SecretPath:         path,
		AttachToProcessEnv: false,
	})
}

func (c *Client) UploadEnv(projectID, env, path string, envMap map[string]string) error {
	existing, err := c.ListSecrets(projectID, env, path)
	if err != nil {
		return err
	}

	for _, s := range existing {
		_, _ = c.sdk.Secrets().Delete(infisical.DeleteSecretOptions{
			ProjectID:   projectID,
			Environment: env,
			SecretPath:  path,
			SecretKey:   s.SecretKey,
		})
	}

	var batch []infisical.BatchCreateSecret
	for k, v := range envMap {
		batch = append(batch, infisical.BatchCreateSecret{
			SecretKey:   k,
			SecretValue: v,
		})
	}

	_, err = c.sdk.Secrets().Batch().Create(infisical.BatchCreateSecretsOptions{
		ProjectID:   projectID,
		Environment: env,
		SecretPath:  path,
		Secrets:     batch,
	})
	return err
}
