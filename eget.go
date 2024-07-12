package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type egetConfig struct {
	arch    string
	os      string
	url     string
	version string
}

func newDefaultConfig() egetConfig {
	url := "https://github.com/zyedidia/eget/releases/download/v%s/eget-%s-%s_%s.tar.gz"
	return egetConfig{arch: runtime.GOARCH, os: runtime.GOOS, url: url, version: "1.3.4"}
}

// TODO: check if dir exists
func createDir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logger.Debug(fmt.Sprintf("creating directory in %s", path))
		if err := os.MkdirAll(path, 0777); err != nil {
			return fmt.Errorf("could not create directory %s", path)
		}
	}
	return nil
}

func downloadEgetBinary(dir string, c egetConfig) error {
	arch := c.arch
	operatingSystem := c.os
	version := c.version
	url := fmt.Sprintf(c.url, version, version, operatingSystem, arch)
	logger.Debug("downloading eget binary", "arch", arch, "os", operatingSystem, "version", version, "url", url)
	tmpDir := path.Join(dir, "tmp")
	if err := createDir(tmpDir); err != nil {
		return err
	}
	egetTar := path.Join(tmpDir, "eget.tar.gz")
	err := downloadFile(url, egetTar)
	if err != nil {
		return err
	}
	err = extractTarGz(egetTar, tmpDir, 1)
	if err != nil {
		return err
	}
	err = makeExecutable(path.Join(tmpDir, "eget"))
	if err != nil {
		return err
	}
	err = rename(path.Join(tmpDir, "eget"), path.Join(dir, "eget"))
	if err != nil {
		return err
	}
	err = os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(url, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download of %s failed with status %d and error %s", url, resp.StatusCode, err)
	}
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func stripLeadingComponents(path string, components int) string {
	parts := strings.Split(path, "/")
	if components >= len(parts) {
		return ""
	}
	return filepath.Join(parts[components:]...)
}

// extractTarGz extracts a tar.gz archive to a destination directory with stripping leading path components.
func extractTarGz(filePath, destPath string, stripComponents int) error {
	// Open the tar.gz file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate through the files in the tar archive
	for {
		header, err := tarReader.Next()

		// If no more files are found, break out of the loop
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Strip leading components from the file path
		strippedPath := stripLeadingComponents(header.Name, stripComponents)
		if strippedPath == "" {
			continue
		}
		path := filepath.Join(destPath, strippedPath)

		// Check the type of the file
		switch header.Typeflag {
		case tar.TypeDir: // if it's a directory
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case tar.TypeReg: // if it's a file
			// Create all directories for the file
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			// Create the file
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			// Copy the file data from the tar archive
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}

func makeExecutable(filePath string) error {
	return os.Chmod(filePath, 0755)
}

func rename(source, dest string) error {
	return os.Rename(source, dest)
}
