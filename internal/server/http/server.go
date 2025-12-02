package internalhttp

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adettelle/image-previewer/config"
	"github.com/adettelle/image-previewer/internal/previewservice"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	cfg  *config.Config
	logg *zap.Logger
	srv  *http.Server
}

type ImageHandler struct {
	PreviewServise Previewer
	CacheCapacity  int
	ScaleOrCrop    string
}

type Previewer interface {
	GeneratePreview(outWidth int, outHeight int,
		imageAddr string, scaleOrCrop string) (previewservice.ResizedImage, error)
	// DownloadFile(filePath string, url string) error
}

func NewServer(cfg *config.Config, logg *zap.Logger) *Server {
	cacheCapacity, err := strconv.Atoi(cfg.CacheCapacity)
	if err != nil {
		log.Fatal(err) // TODO
	}

	ds := previewservice.DownloadService{Logg: logg}

	ps := previewservice.New(cacheCapacity, cfg.PathToSaveIncommingImages,
		cfg.PathToOriginalFile, &ds, logg)

	imageHandler := ImageHandler{
		PreviewServise: ps,
		CacheCapacity:  cacheCapacity,
		ScaleOrCrop:    cfg.Resize,
	}
	router := NewRouter(&imageHandler)
	addr := "0.0.0.0:" + cfg.Port // TODO
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return &Server{cfg: cfg, logg: logg, srv: srv}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("starting http server at", zap.String("address", s.srv.Addr))

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logg.Fatal("server failed: %v", zap.Any("err", err))
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	s.logg.Info("Gracefully shutting down server")
	return nil
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello")) //nolint
}

func (ih *ImageHandler) preview(w http.ResponseWriter, r *http.Request) {
	outWidth := r.PathValue("width")
	outHeight := r.PathValue("height")
	// imageAddr := "https://" + chi.URLParam(r, "*")
	imageAddr := chi.URLParam(r, "*")

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

	resizedImage, err := ih.PreviewServise.GeneratePreview(outW, outH, imageAddr, ih.ScaleOrCrop)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	w.Header().Set("Content-Type", "image/jpeg")              // TODO
	http.ServeFile(w, r, resizedImage.Path+resizedImage.Name) // выдаём наружу // resizedImage
	w.WriteHeader(http.StatusOK)
}
