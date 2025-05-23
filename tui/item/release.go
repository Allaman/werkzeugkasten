package item

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Release struct {
	Tag         string
	PublishedAt time.Time
}

// Mandatory for list.Model
func (r Release) Title() string       { return r.Tag }
func (r Release) Description() string { return r.PublishedAt.Format("2006-01-02") }

func NewRelease(tagName string, publishedAt time.Time) Release {
	return Release{
		Tag:         tagName,
		PublishedAt: publishedAt,
	}
}

func (r Release) FilterValue() string {
	return r.Tag
}

type FetchRelease struct {
	Name        string    `json:"name"`
	TagName     string    `json:"tag_name"`
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	URL         string    `json:"url"`
}

// FetchReleases fetches releases for a GitHub repository
// If token is an empty string, the request will be made without authentication
func FetchReleases(identifier, token string) ([]FetchRelease, error) {
	// NOTE: 100 is max per_page. For more, pagination is needed.
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page=100", identifier)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	// Add authorization header only if token is provided
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			if err == nil {
				err = fmt.Errorf("error closing response body: %w", closeErr)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var releases []FetchRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return releases, nil
}
