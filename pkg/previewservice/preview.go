package previewservice

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"

	"github.com/adettelle/image-previewer/pkg/lru"
)

type PreviewService struct {
	Cache                     lru.LruCache
	PathToSaveIncommingImages string // const "./images/"
	PathToOriginalFile        string // const "/tmp/"
}

func New(cap int, path string, pathToOriginalFile string) *PreviewService {
	return &PreviewService{
		Cache:                     *lru.NewCache(cap),
		PathToSaveIncommingImages: path,
		PathToOriginalFile:        pathToOriginalFile,
	}
}

type ResizedImage struct {
	Path string // const "./images/"
	Name string
}

// returns pathToResizedImage (path + name)
func (ps *PreviewService) GeneratePreview(outWidth int,
	outHeight int, imageAddr string) (ResizedImage, error) {

	originalImageName := base64.StdEncoding.EncodeToString([]byte(imageAddr))
	resizedImageName := originalImageName + "_" + strconv.Itoa(outWidth) + "_" + strconv.Itoa(outHeight)

	resizedImage := ResizedImage{
		Path: ps.PathToSaveIncommingImages,
		Name: resizedImageName,
	}
	// example: ./images/xxx_300_200
	pathToResizedImage := ps.PathToSaveIncommingImages + resizedImageName

	_, ok := ps.Cache.Get(lru.Key(resizedImageName))
	if ok {
		// ps.PathToSaveIncommingImages = pathToResizedImage
		log.Println("returning from cache without putting")
		return resizedImage, nil
	}

	pathToOriginalFile := ps.PathToOriginalFile + originalImageName // "/tmp/"

	err := DownloadFile(pathToOriginalFile, imageAddr)
	if err != nil {
		log.Println("Error downloading file: ", err)
		return ResizedImage{}, err
	}
	log.Println("Downloaded: " + imageAddr)

	err = Scale(pathToOriginalFile, pathToResizedImage, outWidth, outHeight)
	if err != nil && !errors.Is(err, &ResizeError{}) { // TODO CHECK
		log.Println("Error in scaling image: ", err)
		return ResizedImage{}, err
	}

	ps.Cache.Set(lru.Key(resizedImageName), true)

	return resizedImage, nil
}

// сравнивает размер картинки с размером указанным в названии
// func countActualSize() {

// }
