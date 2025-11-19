package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/adettelle/image-previewer/config"
	internalhttp "github.com/adettelle/image-previewer/internal/server/http"
	"github.com/adettelle/image-previewer/pkg/lru"
)

func main() {
	ctx := context.Background()

	cfg := config.New(&ctx)

	lruCapacity, err := strconv.Atoi(cfg.CacheCapacity)
	if err != nil {
		log.Fatal(err)
	}
	ih := internalhttp.ImageHandler{Storager: *lru.NewCache(lruCapacity)}

	router := internalhttp.NewRouter(&ih)

	addr := ":" + cfg.Port // ":8080"
	log.Printf("starting http server at %s\n", addr)

	http.ListenAndServe(addr, router)
}
