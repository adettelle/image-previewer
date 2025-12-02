package main

import (
	"context"
	"log"
	"os"
	"os/signal"
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

	err := os.RemoveAll(cfg.PathToSaveIncommingImages)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(cfg.PathToSaveIncommingImages, 0766)
	if err != nil {
		log.Fatal(err)
	}

	logg := internallogger.GetLogger(cfg.Logger.Level)

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
