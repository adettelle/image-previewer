package integrationtests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/static/x1.jpg"

	resp, err := http.Get(imageAddr)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestDownloadForbidden(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/forbidden/image.jpg"

	resp, err := http.Get(imageAddr)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestDownloadInexistent(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/static/NOimage.jpg"

	resp, err := http.Get(imageAddr)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestHeaderPass(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/header-check/image.jpg"

	// _, err := http.Get(imageAddr)
	// require.NoError(t, err)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, imageAddr, nil)
	require.NoError(t, err)
	req.Header.Add("x-hhh", "123")
	resp, err := client.Do(req) // http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

// http://localhost:8090/static/sample_pdf.pdf
func TestLoadingUnsupportedFile(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/static/sample_pdf.pdf"

	resp, err := http.Get(imageAddr)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusBadRequest)
}
