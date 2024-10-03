package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/creativeprojects/go-selfupdate"
)

func Update(currentVersion string) error {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("allaman/werkzeugkasten"))
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}
	if !found {
		return fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
	}

	if latest.LessOrEqual(currentVersion) {
		slog.Info("current version is the latest version", "version", currentVersion)
		return nil
	}

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return errors.New("could not locate executable path")
	}
	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}
	slog.Info(fmt.Sprintf("Successfully updated to version %s", latest.Version()))
	return nil
}
