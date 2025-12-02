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
