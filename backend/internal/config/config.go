package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIAddr             string
	JWTSecret           string
	TokenTTL            time.Duration
	RedisURL            string
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
		APIAddr:             apiAddr(),
		JWTSecret:           env("JWT_SECRET", "change-me"),
		TokenTTL:            time.Hour,
		RedisURL:            redisURL(),
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

func redisURL() string {
	if value := os.Getenv("REDIS_URL"); value != "" {
		return value
	}

	host := os.Getenv("REDISHOST")
	if host == "" {
		return ""
	}

	port := env("REDISPORT", "6379")
	redisURL := url.URL{
		Scheme: "redis",
		Host:   net.JoinHostPort(host, port),
		Path:   fmt.Sprintf("/%d", envInt("REDIS_DB", 0)),
	}

	user := os.Getenv("REDISUSER")
	password := os.Getenv("REDISPASSWORD")
	if user != "" && password != "" {
		redisURL.User = url.UserPassword(user, password)
	} else if user != "" {
		redisURL.User = url.User(user)
	} else if password != "" {
		redisURL.User = url.UserPassword("", password)
	}

	return redisURL.String()
}

func apiAddr() string {
	if value := os.Getenv("API_ADDR"); value != "" {
		return value
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
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
