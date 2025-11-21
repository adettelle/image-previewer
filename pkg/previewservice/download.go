package previewservice

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadFile will download file from a given url to a filePath.
// It will write as it downloads (useful for large files).
func DownloadFile(filePath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint

	contentType := resp.Header.Get("Content-Type")
	if strings.ToLower(contentType) != "image/jpeg" {
		return fmt.Errorf("unexpected Content-Type %s", contentType)
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close() //nolint

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
