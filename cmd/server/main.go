package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/adettelle/image-previewer/config"
	internallogger "github.com/adettelle/image-previewer/internal/logger"
	internalhttp "github.com/adettelle/image-previewer/internal/server/http"
	"go.uber.org/zap"
)

func main() {
	err := initialize()
	if err != nil {
		log.Fatal(err)
	}
}

func initialize() error {
	startCtx := context.Background()
	ctx, cancel := signal.NotifyContext(startCtx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cfg := config.New(&startCtx)
	logg := internallogger.GetLogger(cfg.Logger.Level)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				err := cleanUp(cfg.PathToOriginalFile, cfg.CleanPeriod)
				if err != nil {
					log.Fatal(err)
				}
			case <-quit:
				return
			}
		}
	}()

	logg.Info("deleting dir for resized", zap.String("dir", cfg.PathToSaveIncommingImages))
	err := os.RemoveAll(cfg.PathToSaveIncommingImages)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("creating dir for resized", zap.String("dir", cfg.PathToSaveIncommingImages))
	err = os.MkdirAll(cfg.PathToSaveIncommingImages, 0766)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("creating temp dir for originals", zap.String("dir", cfg.PathToOriginalFile))
	err = os.MkdirAll(cfg.PathToOriginalFile, 0766)
	if err != nil {
		log.Fatal(err)
	}

	server := internalhttp.NewServer(cfg, logg)

	go func() {
		s := <-ctx.Done()
		logg.Info("Got termination signal: ", zap.Any("Graceful shutdown", s))

		stopCtx, cancel := context.WithTimeout(startCtx, time.Second*3)
		defer cancel()

		err := os.RemoveAll(cfg.PathToSaveIncommingImages)
		if err != nil {
			log.Fatal(err)
		}
		logg.Info("directory is clean")

		if err = server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server", zap.Error(err))
		}

		if err != nil {
			os.Exit(1)
		}

		quit <- struct{}{}
		<-stopCtx.Done()
		os.Exit(0)
	}()

	logg.Info("previewer is running...")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := server.Start(ctx); err != nil {
			logg.Fatal("failed to start http server", zap.Error(err))
		}
	}()

	wg.Wait()
	return nil
}

// path is dir with original files (/tmp/images/)
func cleanUp(path string, seconds int) error {
	now := time.Now()

	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range dirEntry {
		fileInfo, err := entry.Info()
		if err != nil {
			return err
		}
		filtTime := fileInfo.ModTime()
		fullPath := filepath.Join(path, fileInfo.Name())

		if now.Sub(filtTime) > time.Duration(seconds)*time.Second {
			err := os.Remove(fullPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
