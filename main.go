package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jb843051627/beeyard/internal/cache"
	"github.com/jb843051627/beeyard/internal/handler"
	"github.com/jb843051627/beeyard/internal/service"
	"github.com/jb843051627/beeyard/internal/store"
)

func main() {
	dbPath := os.Getenv("BEEYARD_DB")
	if dbPath == "" {
		dbPath = filepath.Join(os.TempDir(), "beeyard.db")
	}
	st, err := store.NewStore(dbPath)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}
	defer st.Close()
	rc := cache.NewReadingCache()
	svc := service.NewService(st, rc)
	h := handler.NewHandler(svc)
	mux := h.Routes()
	addr := os.Getenv("BEEYARD_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("beeyard listening on %s (db=%s)", addr, dbPath)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
