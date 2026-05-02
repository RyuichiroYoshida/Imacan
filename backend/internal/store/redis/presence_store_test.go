package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RyuichiroYoshida/imacan/backend/internal/generated"
	"github.com/RyuichiroYoshida/imacan/backend/internal/presence"
)

func TestPresenceStoreWithRedisAppliesTTLAndDeletesOut(t *testing.T) {
	if os.Getenv("IMACAN_REDIS_INTEGRATION") != "1" {
		t.Skip("set IMACAN_REDIS_INTEGRATION=1 to run Redis integration test")
	}

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	ctx := context.Background()
	client := NewClient(addr, os.Getenv("REDIS_PASSWORD"), 0)
	t.Cleanup(func() {
		_ = client.Close()
	})

	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("ping redis: %v", err)
	}

	store := NewPresenceStore(client)
	service := presence.NewService(store, 105*time.Minute, 20, 30)
	userID := "integration-presence-ttl"
	key := presenceKey(userID)
	t.Cleanup(func() {
		_ = client.Del(ctx, key).Err()
	})
	_ = client.Del(ctx, key).Err()

	classNow := time.Date(2026, 4, 26, 19, 30, 0, 0, time.UTC)
	if _, _, err := service.Update(ctx, userID, generated.CLASS, classNow); err != nil {
		t.Fatal(err)
	}
	assertRedisTTLBetween(t, client.TTL(ctx, key).Val(), 59*time.Minute, 61*time.Minute)

	selfStudyNow := time.Date(2026, 4, 26, 9, 15, 0, 0, time.UTC)
	if _, _, err := service.Update(ctx, userID, generated.SELFSTUDY, selfStudyNow); err != nil {
		t.Fatal(err)
	}
	assertRedisTTLBetween(t, client.TTL(ctx, key).Val(), 11*time.Hour+14*time.Minute, 11*time.Hour+16*time.Minute)

	if _, _, err := service.Update(ctx, userID, generated.OUT, selfStudyNow.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if exists != 0 {
		t.Fatalf("expected OUT to delete presence key %q, exists=%d", key, exists)
	}
}

func assertRedisTTLBetween(t *testing.T, got, min, max time.Duration) {
	t.Helper()

	if got < min || got > max {
		t.Fatalf("expected Redis TTL between %s and %s, got %s", min, max, got)
	}
}
