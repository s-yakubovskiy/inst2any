package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	maxRetries          = 3
	delayBetweenRetries = time.Second * 5
)

func DownloadFile(url string) (io.Reader, error) {
	var resp *http.Response
	var err error
	for i := 0; i < maxRetries; i++ {
		resp, err = http.Get(url)
		if err == nil {
			return resp.Body, nil
		}
		// If error occurs, wait for a while before trying again
		time.Sleep(delayBetweenRetries)
	}
	return nil, fmt.Errorf("failed to download file after %d attempts, last error: %s", maxRetries, err)
}

func DownloadFileToTmp(url string) (string, error) {
	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {
		resp, err = http.Get(url)
		if err == nil {
			break
		}
		time.Sleep(delayBetweenRetries)
	}
	if err != nil {
		return "", fmt.Errorf("failed to download file after %d attempts, last error: %s", maxRetries, err)
	}
	defer resp.Body.Close()

	// Extract filename from URL
	filename := path.Base(url)
	filePath := filepath.Join("/tmp", filename)

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
