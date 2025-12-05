package previewservice

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/adettelle/image-previewer/pkg/lru"
	"go.uber.org/zap"
)

type PreviewService struct {
	Cache               lru.LruCache
	PathToResizedImages string // const "./images/"
	PathToOriginalFile  string // const "/tmp/images/"
	Downloader          Downloader
	Logg                *zap.Logger
}

func New(cap int, pathToResizedImage string, pathToOriginalFile string,
	downloader Downloader, logg *zap.Logger) *PreviewService {

	return &PreviewService{
		Cache:               *lru.NewCache(cap),
		PathToResizedImages: pathToResizedImage,
		PathToOriginalFile:  pathToOriginalFile,
		Downloader:          downloader,
		Logg:                logg,
	}
}

type ResizedImage struct {
	Path string // const "./images/"
	Name string
}

type Downloader interface {
	DownloadFile(filePath string, url string, headers http.Header) error
}

type DownloadService struct {
	Logg   *zap.Logger
	Client http.Client
}

func NewDownloadService(logg *zap.Logger) *DownloadService {
	return &DownloadService{
		Logg:   logg,
		Client: http.Client{},
	}
}

// returns pathToResizedImage (path + name)
// imageAddr is like "https://" + chi.URLParam(r, "*")
func (ps *PreviewService) GeneratePreview(outWidth int,
	outHeight int, imageAddr string, scaleOrCrop string, headers http.Header) (ResizedImage, error) {

	originalImageName := base64.StdEncoding.EncodeToString([]byte(imageAddr))
	resizedImageName := originalImageName + "_" + strconv.Itoa(outWidth) + "_" + strconv.Itoa(outHeight)

	resizedImage := ResizedImage{
		Path: ps.PathToResizedImages, // ./images/
		Name: resizedImageName,       // xxx_300_200
	}
	// example of pathToResizedImage: ./images/xxx_300_200
	pathToResizedImage := ps.PathToResizedImages + resizedImageName

	_, ok := ps.Cache.Get(lru.Key(resizedImageName))
	if ok {
		ps.Logg.Error("returning from cache without putting", zap.String("resizedImageName", resizedImageName))
		return resizedImage, nil
	}

	// example of pathToOriginalFile: "/tmp/images/" + "xxx_300_200"
	pathToOriginalFile := ps.PathToOriginalFile + originalImageName

	fileInfo, err := os.Stat(pathToOriginalFile)
	if os.IsNotExist(err) {
		ps.Logg.Info("file does not exists, should be downloaded:",
			zap.String("pathToOriginalFile", pathToOriginalFile))

		err := ps.Downloader.DownloadFile(pathToOriginalFile, imageAddr, headers)
		if err != nil {
			ps.Logg.Error("error downloading file: ", zap.String("url", imageAddr), zap.Error(err))
			return ResizedImage{}, err
		}
		ps.Logg.Info("Downloaded: ", zap.String("url", imageAddr), zap.String("in", pathToOriginalFile))

	} else if err != nil {
		ps.Logg.Error("error in getting info: ", zap.String("file", pathToOriginalFile), zap.Error(err))
		return ResizedImage{}, err
	} else {
		ps.Logg.Info("file exists, should not be downloaded:",
			zap.String("pathToOriginalFile", pathToOriginalFile),
			zap.String("fileName", fileInfo.Name()))
	}

	switch scaleOrCrop {
	case "scale":
		err = ps.scale(pathToOriginalFile, pathToResizedImage, outWidth, outHeight)
		if err != nil && !errors.Is(err, ResizeError{}) {
			ps.Logg.Error("error in scaling: ", zap.String("image", pathToOriginalFile), zap.Error(err))
			return ResizedImage{}, err
		}
	case "crop":
		err = ps.crop(pathToOriginalFile, pathToResizedImage, outWidth, outHeight)
		if err != nil && !errors.Is(err, ResizeError{}) {
			ps.Logg.Error("error in cropping: ", zap.String("image", pathToOriginalFile), zap.Error(err))
			return ResizedImage{}, err
		}
	}

	ps.Cache.Set(lru.Key(resizedImageName), true)

	return resizedImage, nil
}

// DownloadFile will download file from a given url to a filePath.
// It will write as it downloads (useful for large files).
func (ds *DownloadService) DownloadFile(filePath string, url string, headers http.Header) error {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ds.Logg.Error("error in creating request: ", zap.String("url", url), zap.Error(err))
		return err
	}
	req.Header = headers
	// Get the data
	resp, err := ds.Client.Do(req)
	if err != nil {
		ds.Logg.Error("error in getting url: ", zap.String("url", url), zap.Error(err))
		return err
	}
	defer resp.Body.Close() //nolint

	contentType := resp.Header.Get("Content-Type")
	if strings.ToLower(contentType) != "image/jpeg" {
		ds.Logg.Error("unexpected Content-Type: ", zap.String("contentType", contentType))
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		ds.Logg.Error(" !!!!!!!!!!! error in creating file: ", zap.String("filePath", filePath), zap.Error(err))
		return err
	}
	defer out.Close() //nolint

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		ds.Logg.Error("error in coping file: ", zap.String("file", out.Name()), zap.Error(err))
		return err
	}

	return nil
}
