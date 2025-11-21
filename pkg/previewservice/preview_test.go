package previewservice

import (
	"os"
	"testing"

	"github.com/adettelle/image-previewer/pkg/lru"
	"github.com/c2fo/testify/require"
)

// // приходит картинка 1024x504
// // желаемый размер уменьшенного изображения: 300_200
// // http://localhost:8080/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
// // в кэш записывается имя картинки name(based64)_300_200,
// // сама она сохраняется в ./images/
// // при этом размер картики не 300х200, но один из размеров 300 или 200
// func TestGetNonexistentImageFromCache(t *testing.T) {
// 	ps := New(5)

// 	imageAddr := "https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/examples/image-previewer/gopher_2000x1000.jpg"

// 	pathToResizedImage, err := ps.GeneratePreview(300, 200, imageAddr)
// 	require.NoError(t, err)

// 	_, ok := ps.Cache.Get(lru.Key(pathToResizedImage))
// 	require.True(t, ok)
// }

// приходит картинка 1024x504
// желаемый размер уменьшенного изображения: 300_200
// http://localhost:8080/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
// в кэш записывается имя картинки name(based64)_300_200,
// сама она сохраняется в ./images/
// при этом размер картики не 300х200, но один из размеров 300 или 200
func TestSaveNewIncomingImageToCache(t *testing.T) {
	pathResized := "/tmp/images1/"
	pathToOriginalFile := "/tmp/imagesOriginal1/"

	ps := New(5, pathResized, pathToOriginalFile)

	err := os.MkdirAll(pathResized, 0733)
	require.NoError(t, err)

	err = os.MkdirAll(pathToOriginalFile, 0733)
	require.NoError(t, err)

	imageAddr := "https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/examples/image-previewer/gopher_2000x1000.jpg"

	resizedImage, err := ps.GeneratePreview(444, 222, imageAddr)
	require.NoError(t, err)

	_, ok := ps.Cache.Get(lru.Key(resizedImage.Name))
	require.True(t, ok)

	err = os.RemoveAll(pathResized)
	require.NoError(t, err)
	err = os.RemoveAll(pathToOriginalFile)
	require.NoError(t, err)
}

// приходит картинка 1024x504
// желаемый размер уменьшенного изображения: 300_200
// http://localhost:8080/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
// в кэше такая картинка с именем name(based64)_300_200 есть,
// сама она сохраняется в ./images/
// при этом размер картики не 300х200, но один из размеров 300 или 200
func TestGetIncomingImageFromCache(t *testing.T) {

}

// приходит картинка 256x126
// желаемый размер уменьшенного изображения: 300_200
// http://localhost:8080/fill/300/200/https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/examples/image-previewer/gopher_256x126.jpg
// в кэш записывается имя картинки name(based64)_300_200,
// сама она сохраняется в ./images/
// при этом размер картики не 300х200, а исходный 256x126
func TestSaveNewIncomingImageToCacheWithoutResize(t *testing.T) {

}

// приходит картинка очень большого размера ??????
// желаемый размер уменьшенного изображения: 300_200
//
// она не сохраняется в ./images/
// возвращается ошибка и что еще?
func TestTooBigIncomingImage(t *testing.T) {

}

// надо ли проверять приход png / txt / video ???
