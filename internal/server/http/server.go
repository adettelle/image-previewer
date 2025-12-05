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

// Server represents an HTTP server used to expose the image preview functionality.
type Server struct {
	cfg  *config.Config
	logg *zap.Logger
	srv  *http.Server
}

// ImageHandler handles HTTP requests related to image preview generation.
// It delegates the processing to a Previewer implementation and stores
// configuration options used for resizing.
type ImageHandler struct {
	PreviewServise Previewer
	CacheCapacity  int
	ScaleOrCrop    string
}

// Previewer defines the interface for generating resized image previews.
// Implementations must return a resized image object based on the requested
// dimensions, image source address, and additional processing options.
type Previewer interface {
	GeneratePreview(outWidth int, outHeight int,
		imageAddr string, scaleOrCrop string,
		headers http.Header) (previewservice.ResizedImage, error)
}

// NewServer creates and configures a new Server instance using the provided
// application configuration and logger. It initializes the preview service,
// HTTP router, and the underlying HTTP server.
func NewServer(cfg *config.Config, logg *zap.Logger) *Server {
	cacheCapacity, err := strconv.Atoi(cfg.CacheCapacity)
	if err != nil {
		log.Fatal(err)
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
	addr := "0.0.0.0:" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return &Server{cfg: cfg, logg: logg, srv: srv}
}

// Start launches the HTTP server and begins listening for incoming requests.
// The method blocks until the provided context is canceled or the server fails.
func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("starting http server at", zap.String("address", s.srv.Addr))

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logg.Fatal("server failed: %v", zap.Any("err", err))
	}
	<-ctx.Done()
	return nil
}

// Stop gracefully shuts down the HTTP server using the given context.
// It ensures that all active connections are properly terminated.
func (s *Server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	s.logg.Info("Gracefully shutting down server")
	return nil
}

func mainPage(w http.ResponseWriter, _ *http.Request) {
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

	// -------------
	httpHeaders := http.Header{}
	for key, values := range r.Header {
		for _, value := range values {
			httpHeaders.Add(key, value)
		}
	}

	resizedImage, err := ih.PreviewServise.GeneratePreview(outW, outH, imageAddr, ih.ScaleOrCrop, httpHeaders)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")              // TODO
	http.ServeFile(w, r, resizedImage.Path+resizedImage.Name) // выдаём наружу // resizedImage
	w.WriteHeader(http.StatusOK)
}
