package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/adettelle/image-previewer/config"
	internallogger "github.com/adettelle/image-previewer/internal/logger"
	"github.com/adettelle/image-previewer/internal/previewservice"
	internalhttp "github.com/adettelle/image-previewer/internal/server/http"
	"go.uber.org/zap"
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

	addr := ":" + cfg.Port
	logg := internallogger.GetLogger(cfg.Logger.Level)
	logg.Info("starting http server at", zap.String("port", addr))

	ds := previewservice.DownloadService{}
	ps := previewservice.New(cacheCapacity, pathToSaveIncommingImages,
		pathToOriginalFile, &ds, logg)

	ih := internalhttp.ImageHandler{
		PreviewServise: ps,
		CacheCapacity:  cacheCapacity,
		ScaleOrCrop:    cfg.Resize,
	}

	router := internalhttp.NewRouter(&ih)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
