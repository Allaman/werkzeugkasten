package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func FetchReadme(urlTemplate string) (string, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	branches := []string{"main", "master"}
	filenames := []string{"README.md", "readme.md", "Readme.md"}

	for _, branch := range branches {
		for _, filename := range filenames {
			url := fmt.Sprintf(urlTemplate, branch, filename)
			body, ok, err := fetchURL(client, url)
			if err != nil {
				return "", err
			}
			if ok {
				return body, nil
			}
		}
	}

	return "", fmt.Errorf("could not fetch README")
}

func fetchURL(client http.Client, url string) (string, bool, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return "", false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("status code for downloading README was not ok: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}
	return string(body), true, nil
}
