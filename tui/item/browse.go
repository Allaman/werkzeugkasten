package item

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// openInBrowser opens the given tool identifier in the system's default browser
func OpenInBrowser(identifier string) error {
	url := constructURL(identifier)
	if url == "" {
		return fmt.Errorf("unsupported identifier format: %s", identifier)
	}
	return openURL(url)
}

// constructURL converts a tool identifier to a proper URL
func constructURL(identifier string) string {
	if strings.HasPrefix(identifier, "https://") || strings.HasPrefix(identifier, "http://") {
		return identifier
	}
	if isGitHubRepo(identifier) {
		return fmt.Sprintf("https://github.com/%s", identifier)
	}
	return ""
}

// isGitHubRepo checks if the identifier looks like a GitHub repository identifier (owner/repo)
func isGitHubRepo(identifier string) bool {
	parts := strings.Split(identifier, "/")
	if len(parts) != 2 {
		return false
	}
	for _, part := range parts {
		if part == "" || strings.Contains(part, " ") {
			return false
		}
	}
	return true
}

// openURL opens the given URL using the system's default browser
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "start"
		args = append(args, "")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if _, err := exec.LookPath(cmd); err != nil {
		return fmt.Errorf("no browser opener found (%s): %v", cmd, err)
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
