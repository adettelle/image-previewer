package previewservice

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/adettelle/image-previewer/pkg/lru"
)

type PreviewService struct {
	Cache                     lru.LruCache
	PathToSaveIncommingImages string // const "./images/"
	PathToOriginalFile        string // const "/tmp/"
	Downloader                Downloader
}

func New(cap int, pathToResizedImage string, pathToOriginalFile string, downloader Downloader) *PreviewService {
	return &PreviewService{
		Cache:                     *lru.NewCache(cap),
		PathToSaveIncommingImages: pathToResizedImage,
		PathToOriginalFile:        pathToOriginalFile,
		Downloader:                downloader,
	}
}

type ResizedImage struct {
	Path string // const "./images/"
	Name string
}

type Downloader interface {
	DownloadFile(filePath string, url string) error
}

type DownloadService struct {
}

// returns pathToResizedImage (path + name)
func (ps *PreviewService) GeneratePreview(outWidth int,
	outHeight int, imageAddr string, scaleOrCrop string) (ResizedImage, error) {

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

	pathToOriginalFile := ps.PathToOriginalFile + originalImageName
	// "/tmp/" + "xxx_300_200"

	err := ps.Downloader.DownloadFile(pathToOriginalFile, imageAddr)
	if err != nil {
		log.Println("Error downloading file: ", err)
		return ResizedImage{}, err
	}
	log.Println("Downloaded: " + imageAddr)

	switch scaleOrCrop {
	case "scale":
		err = Scale(pathToOriginalFile, pathToResizedImage, outWidth, outHeight)
		if err != nil && !errors.Is(err, ResizeError{}) { // TODO CHECK
			log.Println("Error in scaling image: ", err)
			return ResizedImage{}, err
		}
	case "crop":
		err = Crop(pathToOriginalFile, pathToResizedImage, outWidth, outHeight)
		if err != nil && !errors.Is(err, ResizeError{}) { // TODO CHECK
			log.Println("Error in cropping image: ", err)
			return ResizedImage{}, err
		}
	}

	ps.Cache.Set(lru.Key(resizedImageName), true)

	return resizedImage, nil
}

// DownloadFile will download file from a given url to a filePath.
// It will write as it downloads (useful for large files).
func (ds *DownloadService) DownloadFile(filePath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint

	contentType := resp.Header.Get("Content-Type")
	if strings.ToLower(contentType) != "image/jpeg" {
		return fmt.Errorf("unexpected Content-Type %s", contentType)
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close() //nolint

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
