package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/RyuichiroYoshida/imacan/backend/internal/auth"
	"github.com/RyuichiroYoshida/imacan/backend/internal/config"
	"github.com/RyuichiroYoshida/imacan/backend/internal/presence"
	"github.com/RyuichiroYoshida/imacan/backend/internal/server"
	redistore "github.com/RyuichiroYoshida/imacan/backend/internal/store/redis"
)

func main() {
	cfg := config.Load()

	authService := auth.NewService(cfg.JWTSecret, cfg.TokenTTL)
	authService.ConfigureDiscord(auth.DiscordConfig{
		ClientID:     cfg.DiscordClientID,
		ClientSecret: cfg.DiscordClientSecret,
		RedirectURI:  cfg.DiscordRedirectURI,
		TokenURL:     cfg.DiscordTokenURL,
		UserURL:      cfg.DiscordUserURL,
	}, nil)

	redisClient := redistore.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	redisDescription := cfg.RedisAddr
	if cfg.RedisURL != "" {
		var err error
		redisClient, err = redistore.NewClientFromURL(cfg.RedisURL)
		if err != nil {
			log.Fatalf("parse redis url: %v", err)
		}
		redisDescription = "REDIS_URL"
	}
	defer redisClient.Close()

	pingCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		log.Fatalf("connect redis %s: %v", redisDescription, err)
	}

	presenceStore := redistore.NewPresenceStore(redisClient)
	presenceService := presence.NewService(presenceStore, cfg.ClassTTL, cfg.SchoolCloseHour, cfg.SchoolCloseMinute)
	handler := server.NewHandler(authService, presenceService)
	router := server.NewRouter(handler, authService)

	log.Printf("ProjectImacan API listening on %s", cfg.APIAddr)
	if err := http.ListenAndServe(cfg.APIAddr, router); err != nil {
		log.Fatal(err)
	}
}
