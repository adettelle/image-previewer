package internalhttp

import (
	"log"
	"net/http"
	"strconv"

	"github.com/adettelle/image-previewer/pkg/previewservice"
	"github.com/go-chi/chi/v5"
)

type ImageHandler struct {
	PreviewServise *previewservice.PreviewService
	CacheCapacity  int
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("start page")
	w.Write([]byte("Hello")) //nolint
}

func (ih *ImageHandler) preview(w http.ResponseWriter, r *http.Request) {
	outWidth := r.PathValue("width")
	outHeight := r.PathValue("height")
	imageAddr := "https://" + chi.URLParam(r, "*")

	outW, err := strconv.Atoi(outWidth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	outH, err := strconv.Atoi(outHeight)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resizedImage, err := ih.PreviewServise.GeneratePreview(outW, outH, imageAddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			// fmt.Printf("Print %s: %s\n", key, value)
			w.Header().Set(key, value)
		}
	}

	w.Header().Set("Content-Type", "image/jpeg")              // TODO
	http.ServeFile(w, r, resizedImage.Path+resizedImage.Name) // выдаём наружу // resizedImage
	w.WriteHeader(http.StatusOK)
}
