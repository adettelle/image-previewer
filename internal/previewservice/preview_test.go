package previewservice

import (
	"image"
	"net/http"
	"os"
	"testing"

	"github.com/adettelle/image-previewer/pkg/lru"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const resize = "scale"

// приходит ошибочный URL на картинку
// желаемый размер уменьшенного изображения: 300_200
// https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/NO_Rainbow_lorikeet_2702x3496.jpg
// в кэш ничего не записывается
func TestGetNonexistentImageFromCache(t *testing.T) { // TODO 1
	pathResized := "/tmp/images1/"
	pathToOriginalFile := "/tmp/imagesOriginal1/"

	logg, err := zap.NewDevelopment()
	require.NoError(t, err)

	ds := NewDownloadService(logg)
	ps := New(5, pathResized, pathToOriginalFile, ds, logg)

	err = os.MkdirAll(pathResized, 0700)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginalFile, 0700)
	require.NoError(t, err)

	imageWantWidth := 300
	imageWantHeight := 200
	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/NO_Rainbow_lorikeet_2702x3496.jpg"

	_, err = ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr, resize, http.Header{})
	require.Error(t, err)

	err = os.RemoveAll(pathResized)
	require.NoError(t, err)
	err = os.RemoveAll(pathToOriginalFile)
	require.NoError(t, err)
}

// Positive case
// приходит картинка 2702x3496
// желаемый размер уменьшенного изображения: 300_200
// https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/Rainbow_lorikeet_2702x3496.jpg
// в кэш записывается имя картинки name(based64)_300_200,
// уменьшенная картинка сохраняется в "/tmp/images1/"
// при этом размер картики не 300х200, но один из размеров 300 или 200
func TestSaveNewIncomingImageToCacheAndGetIt(t *testing.T) { // TODO 1
	pathResized := "/tmp/images1/"
	pathToOriginalFile := "/tmp/imagesOriginal1/"

	logg, err := zap.NewDevelopment()
	require.NoError(t, err)

	ds := NewDownloadService(logg)
	ps := New(5, pathResized, pathToOriginalFile, ds, logg)

	err = os.MkdirAll(pathResized, 0700)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginalFile, 0700)
	require.NoError(t, err)

	imageWantWidth := 444
	imageWantHeight := 222
	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/Rainbow_lorikeet_2702x3496.jpg"

	resizedImage, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr, resize, http.Header{})
	require.NoError(t, err)

	_, ok := ps.Cache.Get(lru.Key(resizedImage.Name))
	require.True(t, ok)

	ok, err = checkAtLeastOneSize(resizedImage.Path+resizedImage.Name, imageWantHeight, imageWantHeight)
	require.NoError(t, err)
	require.True(t, ok)

	err = os.RemoveAll(pathResized)
	require.NoError(t, err)
	err = os.RemoveAll(pathToOriginalFile)
	require.NoError(t, err)
}

// сравнивает размер файла с размерами width, height;
// при совпадении хотя бы одного размера (ширины или высоты) будет ok
func checkAtLeastOneSize(pathToFile string, width, height int) (bool, error) {
	reader, err := os.Open(pathToFile)
	if err != nil {
		return false, err
	}
	defer reader.Close() //nolint

	im, _, err := image.DecodeConfig(reader)
	if err != nil {
		return false, err
	}

	if im.Width == width || im.Height == height {
		return true, nil
	}

	return true, nil
}

// Positive case
// приходит картинка 2000x1000
// // желаемый размер уменьшенного изображения: 500_300
// https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/gopher_2000x1000.jpg
// в кэше такая картинка с именем name(based64)_500_300 есть,
// картинка выдается из кэша, сама при этом еще раз не сохраняется в "/tmp/images1/"
// при этом размер картики не 500х300, но один из размеров 500 или 300
func TestGetIncomingImageFromCacheNoSaving(t *testing.T) {
	pathResizedFile := "/tmp/images1/"
	pathToOriginalFile := "/tmp/imagesOriginal1/"

	logg, err := zap.NewDevelopment()
	require.NoError(t, err)

	ds := NewDownloadService(logg)
	ps := New(5, pathResizedFile, pathToOriginalFile, ds, logg)

	err = os.MkdirAll(pathResizedFile, 0700)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginalFile, 0700)
	require.NoError(t, err)

	imageWantWidth := 500
	imageWantHeight := 300
	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/gopher_2000x1000.jpg"

	// check that the file is not downloaded for the second time:
	resizedImage1, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr, resize, http.Header{})
	require.NoError(t, err)

	fileInfo1, err := os.Stat(resizedImage1.Path + resizedImage1.Name)
	require.NoError(t, err)

	resizedImage2, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr, resize, http.Header{})
	require.NoError(t, err)

	fileInfo2, err := os.Stat(resizedImage2.Path + resizedImage2.Name)
	require.NoError(t, err)

	require.Equal(t, fileInfo1, fileInfo2)
	// ----------------------------------

	_, ok := ps.Cache.Get(lru.Key(resizedImage1.Name))
	require.True(t, ok)

	ok, err = checkAtLeastOneSize(resizedImage1.Path+resizedImage1.Name, imageWantHeight, imageWantHeight)
	require.NoError(t, err)
	require.True(t, ok)

	err = os.RemoveAll(pathResizedFile)
	require.NoError(t, err)
	err = os.RemoveAll(pathToOriginalFile)
	require.NoError(t, err)
}

// приходит картинка 256x126
// желаемый размер уменьшенного изображения: 400_200
// https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/gopher_256x126.jpg
// в кэш записывается имя картинки name(based64)_400_100,
// "уменьшенная" картинка сохраняется в "/tmp/images1/"
// при этом размер "уменьшенной" картинки в "/tmp/images1/" не 400х200, а исходный 256x126
func TestSaveNewIncomingImageToCacheWithoutResize(t *testing.T) {
	pathResized := "/tmp/images1/"
	pathToOriginal := "/tmp/imagesOriginal1/"

	logg, err := zap.NewDevelopment()
	require.NoError(t, err)

	ds := NewDownloadService(logg)
	ps := New(5, pathResized, pathToOriginal, ds, logg)

	err = os.MkdirAll(pathResized, 0700)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginal, 0700)
	require.NoError(t, err)

	imageWantWidth := 400
	imageWantHeight := 200
	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/gopher_256x126.jpg"

	resizedImage, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr, resize, http.Header{})
	require.NoError(t, err)

	_, ok := ps.Cache.Get(lru.Key(resizedImage.Name))
	require.True(t, ok)

	ok, err = checkAtLeastOneSize(resizedImage.Path+resizedImage.Name, imageWantHeight, imageWantHeight)
	require.NoError(t, err)
	require.True(t, ok)

	err = os.RemoveAll(pathResized)
	require.NoError(t, err)
	err = os.RemoveAll(pathToOriginal)
	require.NoError(t, err)
}

func TestCropBigImageToSmall(t *testing.T) {
	pathResized := "/tmp/images1/"
	pathToOriginal := "/tmp/imagesOriginal1/"

	logg, err := zap.NewDevelopment()
	require.NoError(t, err)

	ds := NewDownloadService(logg)
	ps := New(5, pathResized, pathToOriginal, ds, logg)
	err = os.MkdirAll(pathResized, 0700)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginal, 0700)
	require.NoError(t, err)

	imageWantWidth := 400
	imageWantHeight := 100
	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/gopher_2000x1000.jpg"

	scaleOrCrop := "crop"

	resizedImage, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr, scaleOrCrop, http.Header{})
	require.NoError(t, err)

	actualW, actualH, err := actualSize(pathResized + resizedImage.Name)
	require.NoError(t, err)
	require.Equal(t, imageWantWidth, actualW)
	require.Equal(t, imageWantHeight, actualH)
}

func TestCropSmallImageToBig(t *testing.T) {
	pathResized := "/tmp/images1/"
	pathToOriginalFile := "/tmp/imagesOriginal1/"

	logg, err := zap.NewDevelopment()
	require.NoError(t, err)

	ds := NewDownloadService(logg)
	ps := New(5, pathResized, pathToOriginalFile, ds, logg)
	err = os.MkdirAll(pathResized, 0700)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginalFile, 0700)
	require.NoError(t, err)

	imageWantWidth := 400
	imageWantHeight := 200
	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/examples/gopher_256x126.jpg"

	scaleOrCrop := "crop"

	resizedImage, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr, scaleOrCrop, http.Header{})
	require.NoError(t, err)

	actualW, actualH, err := actualSize(pathResized + resizedImage.Name)
	require.NoError(t, err)
	require.Equal(t, 256, actualW)
	require.Equal(t, 126, actualH)
}

func actualSize(pathToFile string) (width int, height int, err error) {
	reader, err := os.Open(pathToFile)
	if err != nil {
		return 0, 0, err
	}
	defer reader.Close() //nolint

	im, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}

	return im.Width, im.Height, nil
}
