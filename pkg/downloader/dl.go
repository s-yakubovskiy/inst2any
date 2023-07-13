package downloader

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxRetries = 3
const delayBetweenRetries = time.Second * 5

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
