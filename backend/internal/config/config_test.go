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
	t.Setenv("REDISHOST", "ignored.railway.internal")

	cfg := Load()

	if cfg.RedisURL != "redis://default:secret@redis.railway.internal:6379/0" {
		t.Fatalf("RedisURL = %q", cfg.RedisURL)
	}
}

func TestLoadBuildsRedisURLFromRailwayRedisVariables(t *testing.T) {
	t.Setenv("REDIS_URL", "")
	t.Setenv("REDISHOST", "redis.railway.internal")
	t.Setenv("REDISPORT", "6379")
	t.Setenv("REDISUSER", "default")
	t.Setenv("REDISPASSWORD", "secret")
	t.Setenv("REDIS_DB", "0")

	cfg := Load()

	if cfg.RedisURL != "redis://default:secret@redis.railway.internal:6379/0" {
		t.Fatalf("RedisURL = %q", cfg.RedisURL)
	}
}

func TestLoadBuildsRedisURLFromRailwayRedisVariablesWithEscaping(t *testing.T) {
	t.Setenv("REDIS_URL", "")
	t.Setenv("REDISHOST", "redis.railway.internal")
	t.Setenv("REDISPORT", "6379")
	t.Setenv("REDISUSER", "default")
	t.Setenv("REDISPASSWORD", "sec/ret")
	t.Setenv("REDIS_DB", "2")

	cfg := Load()

	if cfg.RedisURL != "redis://default:sec%2Fret@redis.railway.internal:6379/2" {
		t.Fatalf("RedisURL = %q", cfg.RedisURL)
	}
}
