package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func FetchReadme(url string) (string, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	var resp *http.Response
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code for downloading README was not ok: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
