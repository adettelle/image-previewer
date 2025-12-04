package previewservice

import (
	"encoding/base64"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDownload(t *testing.T) {
	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/gopher_256x126.jpg"
	path := "/tmp/images2/"
	logg, err := zap.NewDevelopment()
	require.NoError(t, err)

	ds := NewDownloadService(logg)
	ps := New(5, "", path, ds, logg)

	err = os.MkdirAll(path, 0700)
	require.NoError(t, err)

	originalImageName := base64.StdEncoding.EncodeToString([]byte(imageAddr))
	pathToOriginalFile := path + originalImageName

	err = ps.Downloader.DownloadFile(pathToOriginalFile, imageAddr, http.Header{})
	require.NoError(t, err)

	file, err := os.Open(pathToOriginalFile)
	require.NoError(t, err)
	defer file.Close() // nolint

	err = os.RemoveAll(path)
	require.NoError(t, err)
}
