### Usage

**env [dev, staging, prod]**

- envctl list
- envctl push --env {env}
- envctl pull --env {env}

#### Behavior

**push**
- Reads local `.env`
- Sends all key/value pairs to Infisical for the selected environment (`--env`)
- Remote values in Infisical will be **replaced** if the same key already exists

**pull**
- Fetches all secrets from Infisical for the selected environment (`--env`)
- Rewrites local `.env` with the remote values
- The existing `.env` file will be полностью replaced by the latest version from Infisical

#### Notes

- `pull` will overwrite the existing `.env` file with the latest values from Infisical.
- `push` will overwrite existing keys in Infisical if they already exist remotely.
- Both operations are destructive with respect to previous values — make sure you have backups if needed.

