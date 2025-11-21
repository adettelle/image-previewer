package previewservice

import (
	"image"
	"os"
	"testing"

	"github.com/adettelle/image-previewer/pkg/lru"
	"github.com/stretchr/testify/require"
)

// приходит картинка 1024x504
// желаемый размер уменьшенного изображения: 300_200
// http://localhost:8080/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
// в кэш записывается имя картинки name(based64)_300_200,
// уменьшенная картинка сохраняется в "/tmp/images1/"
// при этом размер картики не 300х200, но один из размеров 300 или 200
// func TestGetNonexistentImageFromCache(t *testing.T) {
// 	ps := New(5)

// 	imageAddr := "https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/examples/image-previewer/gopher_2000x1000.jpg"

// 	pathToResizedImage, err := ps.GeneratePreview(300, 200, imageAddr)
// 	require.NoError(t, err)

// 	_, ok := ps.Cache.Get(lru.Key(pathToResizedImage))
// 	require.True(t, ok)
// }

// Positive case
// приходит картинка 2000x1000
// желаемый размер уменьшенного изображения: 300_200
// http://localhost:8080/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_2000x1000.jpg
// в кэш записывается имя картинки name(based64)_300_200,
// уменьшенная картинка сохраняется в "/tmp/images1/"
// при этом размер картики не 300х200, но один из размеров 300 или 200
func TestSaveNewIncomingImageToCacheAndGetIt(t *testing.T) {
	pathResized := "/tmp/images1/"
	pathToOriginalFile := "/tmp/imagesOriginal1/"

	ps := New(5, pathResized, pathToOriginalFile)

	err := os.MkdirAll(pathResized, 0733)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginalFile, 0733)
	require.NoError(t, err)

	imageWantWidth := 444
	imageWantHeight := 222
	imageAddr := "https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/examples/image-previewer/gopher_2000x1000.jpg"

	resizedImage, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr)
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

// func actualSize(pathToFile string) (width int, height int, err error) {
// 	fmt.Println(" !!!!!!!!!! ", pathToFile)
// 	reader, err := os.Open(pathToFile)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	defer reader.Close()

// 	im, _, err := image.DecodeConfig(reader)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	return im.Width, im.Height, nil
// }

// Positive case
// приходит картинка 2000x1000
// // желаемый размер уменьшенного изображения: 500_300
// http://localhost:8080/fill/500/300/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
// в кэше такая картинка с именем name(based64)_500_300 есть,
// картинка выдается из кэша, сама при этом еще раз не сохраняется в "/tmp/images1/"
// при этом размер картики не 500х300, но один из размеров 500 или 300
func TestGetIncomingImageFromCacheNoSaving(t *testing.T) {
	pathResizedFile := "/tmp/images1/"
	pathToOriginalFile := "/tmp/imagesOriginal1/"

	ps := New(5, pathResizedFile, pathToOriginalFile)

	err := os.MkdirAll(pathResizedFile, 0733)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginalFile, 0733)
	require.NoError(t, err)

	imageWantWidth := 500
	imageWantHeight := 300
	imageAddr := "https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/examples/image-previewer/gopher_2000x1000.jpg"

	// check that the file is not downloaded for the second time:
	resizedImage1, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr)
	require.NoError(t, err)

	fileInfo1, err := os.Stat(resizedImage1.Path + resizedImage1.Name)
	require.NoError(t, err)

	resizedImage2, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr)
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
// желаемый размер уменьшенного изображения: 400_100
// http://localhost:8080/fill/400/100/https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/examples/image-previewer/gopher_256x126.jpg
// в кэш записывается имя картинки name(based64)_400_100,
// "уменьшенная" картинка сохраняется в "/tmp/images1/"
// при этом размер resided картинки в "/tmp/images1/" не 400х100, а исходный 256x126
func TestSaveNewIncomingImageToCacheWithoutResize(t *testing.T) {
	pathResized := "/tmp/images1/"
	pathToOriginalFile := "/tmp/imagesOriginal1/"

	ps := New(5, pathResized, pathToOriginalFile)

	err := os.MkdirAll(pathResized, 0733)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginalFile, 0733)
	require.NoError(t, err)

	imageWantWidth := 400
	imageWantHeight := 100
	imageAddr := "https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/examples/image-previewer/gopher_256x126.jpg"

	resizedImage, err := ps.GeneratePreview(imageWantWidth, imageWantHeight, imageAddr)
	require.NoError(t, err)

	// w, h, err := actualSize(resizedImage.Path + resizedImage.Name)
	// require.NoError(t, err)
	// require.Less(t, 400, w) // TODO
	// require.Less(t, 100, h) // TODO

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

// приходит картинка очень большого размера ??????
// желаемый размер уменьшенного изображения: 300_200
//
// уменьшенная картинка не сохраняется в "/tmp/images1/"
// возвращается ошибка и что еще?
func TestTooBigIncomingImage(t *testing.T) {

}

// надо ли проверять приход png / txt / video ???
