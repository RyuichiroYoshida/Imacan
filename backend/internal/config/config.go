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
	RedisAddr           string
	RedisPassword       string
	RedisDB             int
	SchoolCloseHour     int
	SchoolCloseMinute   int
	ClassTTL            time.Duration
	SchoolRadiusMeters  int
	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURI  string
	DiscordTokenURL     string
	DiscordUserURL      string
}

func Load() Config {
	closeHour, closeMinute := parseClock(env("SCHOOL_CLOSE_TIME", "20:30"))

	return Config{
		APIAddr:             env("API_ADDR", ":8080"),
		JWTSecret:           env("JWT_SECRET", "change-me"),
		TokenTTL:            time.Hour,
		RedisAddr:           env("REDIS_ADDR", "localhost:6379"),
		RedisPassword:       os.Getenv("REDIS_PASSWORD"),
		RedisDB:             envInt("REDIS_DB", 0),
		SchoolCloseHour:     closeHour,
		SchoolCloseMinute:   closeMinute,
		ClassTTL:            105 * time.Minute,
		SchoolRadiusMeters:  envInt("SCHOOL_RADIUS_METERS", 300),
		DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		DiscordRedirectURI:  os.Getenv("DISCORD_REDIRECT_URI"),
		DiscordTokenURL:     env("DISCORD_TOKEN_URL", "https://discord.com/api/oauth2/token"),
		DiscordUserURL:      env("DISCORD_USER_URL", "https://discord.com/api/users/@me"),
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
