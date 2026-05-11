package infisical

// Credentials baked in at build time via:
//
//	go build -ldflags "
//	  -X 'mit/platform/internal/envctl/infisical.embeddedClientID=...'
//	  -X 'mit/platform/internal/envctl/infisical.embeddedClientSecret=...'
//	  -X 'mit/platform/internal/envctl/infisical.embeddedSiteURL=https://us.infisical.com'
//	"
//
// In CI these come from GitHub Actions secrets (see .github/workflows/release.yml).
// Source files keep them empty so nothing sensitive is ever committed.
var (
	embeddedClientID     = ""
	embeddedClientSecret = ""
	embeddedSiteURL      = ""
)

const defaultSiteURL = "https://us.infisical.com"
