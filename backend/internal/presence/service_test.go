package presence

import (
	"context"
	"testing"
	"time"

	"github.com/RyuichiroYoshida/imacan/backend/internal/generated"
)

func TestUpdateAppliesClassTTLWithoutPassingSchoolClose(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(store, 105*time.Minute, 20, 30)
	now := time.Date(2026, 4, 26, 19, 30, 0, 0, time.UTC)

	record, hasExpiry, err := service.Update(context.Background(), "user-1", generated.CLASS, now)
	if err != nil {
		t.Fatal(err)
	}

	want := time.Date(2026, 4, 26, 20, 30, 0, 0, time.UTC)
	if !hasExpiry || !record.ExpiresAt.Equal(want) {
		t.Fatalf("expected class presence to expire at school close %s, got hasExpiry=%v expiresAt=%s", want, hasExpiry, record.ExpiresAt)
	}
}

func TestUpdateAppliesSelfStudyTTLUntilSchoolClose(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(store, 105*time.Minute, 20, 30)
	now := time.Date(2026, 4, 26, 9, 15, 0, 0, time.UTC)

	record, hasExpiry, err := service.Update(context.Background(), "user-1", generated.SELFSTUDY, now)
	if err != nil {
		t.Fatal(err)
	}

	want := time.Date(2026, 4, 26, 20, 30, 0, 0, time.UTC)
	if !hasExpiry || !record.ExpiresAt.Equal(want) {
		t.Fatalf("expected self-study presence to expire at school close %s, got hasExpiry=%v expiresAt=%s", want, hasExpiry, record.ExpiresAt)
	}
}

func TestUpdateDeletesPresenceWhenActivityIsOut(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	service := NewService(store, 105*time.Minute, 20, 30)
	now := time.Date(2026, 4, 26, 9, 15, 0, 0, time.UTC)

	if _, _, err := service.Update(ctx, "user-1", generated.SELFSTUDY, now); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := store.Get(ctx, "user-1"); err != nil || !ok {
		t.Fatalf("expected presence before OUT update, ok=%v err=%v", ok, err)
	}

	record, hasExpiry, err := service.Update(ctx, "user-1", generated.OUT, now.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if hasExpiry || !record.ExpiresAt.IsZero() {
		t.Fatalf("expected OUT response without expiry, got hasExpiry=%v record=%+v", hasExpiry, record)
	}
	if _, ok, err := store.Get(ctx, "user-1"); err != nil || ok {
		t.Fatalf("expected presence to be deleted after OUT update, ok=%v err=%v", ok, err)
	}
}

func TestCurrentDeletesExpiredPresence(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	service := NewService(store, 105*time.Minute, 20, 30)
	now := time.Date(2026, 4, 26, 20, 29, 0, 0, time.UTC)

	if _, _, err := service.Update(ctx, "user-1", generated.SELFSTUDY, now); err != nil {
		t.Fatal(err)
	}

	_, active, err := service.Current(ctx, "user-1", time.Date(2026, 4, 26, 20, 31, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if active {
		t.Fatal("expected expired presence to be inactive")
	}
	if _, ok, err := store.Get(ctx, "user-1"); err != nil || ok {
		t.Fatalf("expected expired presence to be deleted, ok=%v err=%v", ok, err)
	}
}
