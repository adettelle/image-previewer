package internalhttp

import (
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"strings"

	"image/jpeg"

	"github.com/adettelle/image-previewer/pkg/lru"
	"github.com/go-chi/chi/v5"
)

type ImageHandler struct {
	Storager lru.LruCache
}

func NewRouter(h *ImageHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", hello)
	r.Get("/fill/{width}/{height}/*", h.preview)

	return r
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (h *ImageHandler) preview(w http.ResponseWriter, r *http.Request) {
	width := r.PathValue("width")
	height := r.PathValue("height")
	imageAddr := "https://" + chi.URLParam(r, "*")
	fmt.Println(width, height, imageAddr)

	originalImageName := base64.StdEncoding.EncodeToString([]byte(imageAddr))
	resizedImageName := originalImageName + "_" + width + "_" + height

	_, ok := h.Storager.Get(lru.Key(resizedImageName))
	if ok {
		fmt.Println("Got from cache: " + imageAddr)
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, "./images/"+resizedImageName)
		w.WriteHeader(http.StatusOK)
		return
	}

	path := "/tmp/" + originalImageName // "./images/saveas.jpg"

	err := DownloadFile(path, imageAddr)
	if err != nil {
		fmt.Println("Error downloading file: ", err)
		return
	}
	fmt.Println("Downloaded: " + imageAddr)

	err = GetImageSize(path, "./images/"+resizedImageName)
	if err != nil {
		fmt.Println("Error in getting size of file: ", err)
		return
	}

	h.Storager.Set(lru.Key(resizedImageName), true)

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeFile(w, r, "./images/"+resizedImageName)
	w.WriteHeader(http.StatusOK)
}

// DownloadFile will dpathownload from a given url to a file. It will
// write as it downloads (useful for large files).
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	if strings.ToLower(contentType) != "image/jpeg" {
		return fmt.Errorf("unexpected Content-Type %s", contentType)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func GetImageSize(path string, resizedImagePath string) error {
	fmt.Println("path", path)
	// Open file.
	inputFile, err := os.Open(path)
	if err != nil {
		fmt.Println(" 555555555555555555 ")
		return err
	}

	// -------------------------------------------------------------------------

	originalImage, err := jpeg.Decode(inputFile)
	if err != nil {
		fmt.Println(" 3333333333333333 ")
		return err
	}
	bounds := originalImage.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// cropSize := image.Rect(0, 0, width/2, height/2)
	cropSize := image.Rect(width/4, height/4, width*3/4, height*3/4)

	// cropSize = cropSize.Add(image.Point{100, 100})

	croppedImage := originalImage.(SubImager).SubImage(cropSize)

	// -------------------------------------------------------------------------

	croppedImageFile, err := os.Create(resizedImagePath)
	if err != nil {
		fmt.Println(" 22222222222222 ")
		return err
	}
	defer croppedImageFile.Close()

	err = jpeg.Encode(croppedImageFile, croppedImage, nil)
	if err != nil {
		fmt.Println(" 1111111111111 ")
		return err
	}

	return nil
}
