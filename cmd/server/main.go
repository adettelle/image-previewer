package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/adettelle/image-previewer/config"
	"github.com/adettelle/image-previewer/internal/previewservice"
	internalhttp "github.com/adettelle/image-previewer/internal/server/http"
)

const (
	pathToSaveIncommingImages = "./images/"
	pathToOriginalFile        = "/tmp/"
)

func main() {
	ctx := context.Background()

	cfg := config.New(&ctx)

	cacheCapacity, err := strconv.Atoi(cfg.CacheCapacity)
	if err != nil {
		log.Fatal(err)
	}

	ds := previewservice.DownloadService{}
	ps := previewservice.New(cacheCapacity, pathToSaveIncommingImages, pathToOriginalFile, &ds)

	ih := internalhttp.ImageHandler{
		PreviewServise: ps,
		CacheCapacity:  cacheCapacity,
		ScaleOrCrop:    cfg.Resize,
	}

	router := internalhttp.NewRouter(&ih)

	addr := ":" + cfg.Port
	log.Printf("starting http server at %s\n", addr)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
