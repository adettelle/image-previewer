package internalhttp

import "testing"

// приходит картинка 1024x504
// желаемый размер уменьшенного изображения: 300_200
// http://localhost:8080/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
// в кэш записывается имя картинки name(based64)_300_200,
// сама она сохраняется в ./images/
// при этом размер картики не 300х200, но один из размеров 300 или 200
func TestSaveNewIncomingImageToCache(t *testing.T) {

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
