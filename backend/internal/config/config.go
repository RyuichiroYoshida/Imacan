package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIAddr             string
	JWTSecret           string
	TokenTTL            time.Duration
	SchoolCloseHour     int
	SchoolCloseMinute   int
	ClassTTL            time.Duration
	SchoolRadiusMeters  int
	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURI  string
}

func Load() Config {
	closeHour, closeMinute := parseClock(env("SCHOOL_CLOSE_TIME", "20:30"))

	return Config{
		APIAddr:             env("API_ADDR", ":8080"),
		JWTSecret:           env("JWT_SECRET", "change-me"),
		TokenTTL:            time.Hour,
		SchoolCloseHour:     closeHour,
		SchoolCloseMinute:   closeMinute,
		ClassTTL:            105 * time.Minute,
		SchoolRadiusMeters:  envInt("SCHOOL_RADIUS_METERS", 300),
		DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		DiscordRedirectURI:  os.Getenv("DISCORD_REDIRECT_URI"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseClock(value string) (int, int) {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 20, 30
	}
	return parsed.Hour(), parsed.Minute()
}
