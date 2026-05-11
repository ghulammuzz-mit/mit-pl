package file

import (
	"fmt"
	"os"
	"strings"

	infisical "github.com/infisical/go-sdk"
	"github.com/joho/godotenv"
)

func Write(file string, secrets []infisical.Secret) error {
	var b strings.Builder

	for _, s := range secrets {
		b.WriteString(fmt.Sprintf("%s=%s\n", s.SecretKey, s.SecretValue))
	}

	return os.WriteFile(file, []byte(b.String()), 0644)
}

func Read(file string) (map[string]string, error) {
	return godotenv.Read(file)
}
