package selfupdate

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// downloadFile returns downloaded file path and it's sha256 hash
// the function can download the file of any size
func downloadFile(tempFolder, url string) (string, error) {
	client := http.DefaultClient
	client.Timeout = config.HTTPTimeout + config.DownloadTimeout

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned %d", resp.StatusCode)
	}

	fileName := prepareFileName(url, resp.Header.Get("Content-Disposition"))
	dstFilePath := tempFolder + "/" + fileName

	f, err := os.OpenFile(dstFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}

	return dstFilePath, nil
}

func prepareFileName(url, contentDisposition string) string {
	fallbackFilename := filepath.Base(url)

	if contentDisposition == "" {
		return fallbackFilename
	}

	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return fallbackFilename
	}
	filename, ok := params["filename"]
	if ok {
		return filename
	}
	return fallbackFilename
}
