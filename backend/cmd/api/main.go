package main

import (
	"log"
	"net/http"

	"github.com/RyuichiroYoshida/imacan/backend/internal/auth"
	"github.com/RyuichiroYoshida/imacan/backend/internal/config"
	"github.com/RyuichiroYoshida/imacan/backend/internal/presence"
	"github.com/RyuichiroYoshida/imacan/backend/internal/server"
)

func main() {
	cfg := config.Load()

	authService := auth.NewService(cfg.JWTSecret, cfg.TokenTTL)
	presenceService := presence.NewService(cfg.ClassTTL, cfg.SchoolCloseHour, cfg.SchoolCloseMinute)
	handler := server.NewHandler(authService, presenceService)
	router := server.NewRouter(handler, authService)

	log.Printf("ProjectImacan API listening on %s", cfg.APIAddr)
	if err := http.ListenAndServe(cfg.APIAddr, router); err != nil {
		log.Fatal(err)
	}
}
