package main

import (
	"8th_pract_go/internal/config"
	"8th_pract_go/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	apiCfg := handlers.NewApi(cfg.CookieName, cfg.EncryptionKey)

	router := chi.NewRouter()
	router.Post("/linear", apiCfg.LinearHandler)
	router.Post("/concurrent", apiCfg.ConcurrentHandler)
	router.Get("/get", apiCfg.GetCookie)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
