package internalhttp

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adettelle/image-previewer/pkg/file"
	"github.com/adettelle/image-previewer/pkg/lru"
	"github.com/go-chi/chi/v5"
)

// TODO CHECK !!!!!!!!!!! пока не используется, проверить, почему
const pathToSaveIncommingImages = "./images/"

type ImageHandler struct {
	Storager lru.LruCache
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (h *ImageHandler) preview(w http.ResponseWriter, r *http.Request) {
	outWidth := r.PathValue("width")
	outHeight := r.PathValue("height")
	imageAddr := "https://" + chi.URLParam(r, "*")
	fmt.Println(outWidth, outHeight, imageAddr)

	originalImageName := base64.StdEncoding.EncodeToString([]byte(imageAddr))
	resizedImageName := originalImageName + "_" + outWidth + "_" + outHeight

	// TODO !!!!!!!!!!
	outW, err := strconv.Atoi(outWidth)
	if err != nil {
		fmt.Println(err) // TODO !!!!!!!!!!
	}
	outH, err := strconv.Atoi(outHeight)
	if err != nil {
		fmt.Println(err) // TODO !!!!!!!!!!
	}

	previewService := file.New()

	pathToSave, err := previewService.GeneratePreview(outW, outH, imageAddr) // x = ./images/xxx_20_10
	if err != nil {
		fmt.Println(err) // TODO !!!!!!!!!!
	}

	// _, ok := h.Storager.Get(lru.Key(resizedImageName))
	// if ok {
	// 	fmt.Println("Got from cache: " + imageAddr)
	w.Header().Set("Content-Type", "image/jpeg")
	// было pathToSaveIncommingImages+resizedImageName // TODO CHECK !!!!!!!!!!!
	http.ServeFile(w, r, pathToSave+resizedImageName) // выдаём наружу
	w.WriteHeader(http.StatusOK)
	// return
	// }
	// TODO !ok !!!!!!!!!

	path := "/tmp/" + originalImageName // "./images/saveas.jpg"

	err = file.DownloadFile(path, imageAddr)
	if err != nil {
		fmt.Println("Error downloading file: ", err)
		return
	}
	fmt.Println("Downloaded: " + imageAddr)

	outWidthInPxs, err := strconv.Atoi(outWidth)
	if err != nil {
		return
	}
	outHeightInPxs, err := strconv.Atoi(outHeight)
	if err != nil {
		return
	}

	// err = resize.Crop(path, "./images/"+resizedImageName, widthInPxs, heightInPxs)
	// if err != nil && !errors.Is(err, &file.ResizeError{}) {
	// 	fmt.Println("Error in cropping image: ", err)
	// 	return
	// }

	err = file.Scale(path, pathToSaveIncommingImages+resizedImageName, outWidthInPxs, outHeightInPxs)
	if err != nil && !errors.Is(err, &file.ResizeError{}) { // TODO CHECK
		fmt.Println("Error in scaling image: ", err)
		return
	}

	h.Storager.Set(lru.Key(resizedImageName), true)

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeFile(w, r, "./images/"+resizedImageName)
	w.WriteHeader(http.StatusOK)
}
