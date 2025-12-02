package previewservice

import (
	"encoding/base64"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCrop(t *testing.T) {
	pathResized := "/tmp/images1/"
	pathToOriginal := "/tmp/imagesOriginals1/"
	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/gopher_2000x1000.jpg"
	outWidth := 400
	outHeight := 200

	logg, err := zap.NewDevelopment()
	require.NoError(t, err)

	ds := DownloadService{Logg: logg}
	ps := New(5, pathResized, pathToOriginal, &ds, logg)

	originalImageName := base64.StdEncoding.EncodeToString([]byte(imageAddr))
	resizedImageName := originalImageName + "_" + strconv.Itoa(outWidth) + "_" + strconv.Itoa(outHeight)

	err = os.MkdirAll(pathResized, 0700)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginal, 0700)
	require.NoError(t, err)

	pathToOriginalFile := pathToOriginal + originalImageName

	err = ps.Downloader.DownloadFile(pathToOriginalFile, imageAddr)
	require.NoError(t, err)

	err = ps.crop(pathToOriginalFile, pathResized+resizedImageName, outWidth, outHeight)
	require.NoError(t, err)

	actualW, actualH, err := actualSize(pathResized + resizedImageName)
	require.NoError(t, err)
	require.Equal(t, outWidth, actualW)
	require.Equal(t, outHeight, actualH)

	err = os.RemoveAll(pathResized)
	require.NoError(t, err)
	err = os.RemoveAll(pathToOriginal)
	require.NoError(t, err)
}
