package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultHostURL = "https://us.infisical.com"

type Profile struct {
	ClientID     string
	ClientSecret string
	HostURL      string
}

func Path() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envctl", "config")
}

func Load(profile string) (*Profile, error) {
	f, err := os.Open(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	profiles := parse(f)
	p, ok := profiles[profile]
	if !ok {
		return nil, nil
	}
	return &p, nil
}

func Save(profile string, p Profile) error {
	path := Path()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	// Read existing profiles to preserve others
	existing := map[string]Profile{}
	if f, err := os.Open(path); err == nil {
		existing = parse(f)
		f.Close()
	}
	existing[profile] = p

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for name, prof := range existing {
		fmt.Fprintf(w, "[%s]\n", name)
		fmt.Fprintf(w, "infisical_client_id = %s\n", prof.ClientID)
		fmt.Fprintf(w, "infisical_client_secret = %s\n", prof.ClientSecret)
		hostURL := prof.HostURL
		if hostURL == "" {
			hostURL = defaultHostURL
		}
		fmt.Fprintf(w, "infisical_host_url = %s\n", hostURL)
		fmt.Fprintln(w)
	}
	return w.Flush()
}

func parse(f *os.File) map[string]Profile {
	profiles := map[string]Profile{}
	scanner := bufio.NewScanner(f)
	current := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			current = line[1 : len(line)-1]
			continue
		}
		if current == "" || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		p := profiles[current]
		switch key {
		case "infisical_client_id":
			p.ClientID = val
		case "infisical_client_secret":
			p.ClientSecret = val
		case "infisical_host_url":
			p.HostURL = val
		}
		profiles[current] = p
	}
	return profiles
}

func DefaultHostURL() string {
	return defaultHostURL
}
