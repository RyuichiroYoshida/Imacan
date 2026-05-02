package config

import "testing"

func TestLoadUsesPortWhenAPIAddrUnset(t *testing.T) {
	t.Setenv("API_ADDR", "")
	t.Setenv("PORT", "4321")

	cfg := Load()

	if cfg.APIAddr != ":4321" {
		t.Fatalf("APIAddr = %q, want %q", cfg.APIAddr, ":4321")
	}
}

func TestLoadPrefersAPIAddrOverPort(t *testing.T) {
	t.Setenv("API_ADDR", ":9090")
	t.Setenv("PORT", "4321")

	cfg := Load()

	if cfg.APIAddr != ":9090" {
		t.Fatalf("APIAddr = %q, want %q", cfg.APIAddr, ":9090")
	}
}

func TestLoadReadsRedisURL(t *testing.T) {
	t.Setenv("REDIS_URL", "redis://default:secret@redis.railway.internal:6379/0")

	cfg := Load()

	if cfg.RedisURL != "redis://default:secret@redis.railway.internal:6379/0" {
		t.Fatalf("RedisURL = %q", cfg.RedisURL)
	}
}
