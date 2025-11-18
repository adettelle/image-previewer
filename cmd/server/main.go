package main

import (
	"log"
	"net/http"

	internalhttp "github.com/adettelle/image-previewer/internal/server/http"
	"github.com/adettelle/image-previewer/pkg/lru"
)

func main() {
	lruCapacity := 10
	ih := internalhttp.ImageHandler{Storager: *lru.NewCache(lruCapacity)}

	router := internalhttp.NewRouter(&ih)

	addr := ":8080"
	log.Printf("starting http server at %s\n", addr)

	http.ListenAndServe(addr, router)
}
