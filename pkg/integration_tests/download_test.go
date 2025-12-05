package integrationtests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/static/x1.jpg"

	client := &http.Client{}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, imageAddr, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDownloadForbidden(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/forbidden/image.jpg"

	client := &http.Client{}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, imageAddr, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDownloadInexistent(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/static/NOimage.jpg"

	client := &http.Client{}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, imageAddr, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHeaderPass(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/header-check/image.jpg"

	client := &http.Client{}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, imageAddr, nil)
	require.NoError(t, err)
	req.Header.Add("x-custom-header", "123")
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

// http://localhost:8090/static/sample_pdf.pdf
func TestLoadingUnsupportedFile(t *testing.T) {
	imageAddr := "http://previewer:8080/fill/300/200/http://nginx:80/static/sample_pdf.pdf"

	client := &http.Client{}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, imageAddr, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
