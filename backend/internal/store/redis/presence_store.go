package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/RyuichiroYoshida/imacan/backend/internal/presence"
	goredis "github.com/redis/go-redis/v9"
)

const presenceKeyPrefix = "imacan:presence:"

type PresenceStore struct {
	client *goredis.Client
}

func NewPresenceStore(client *goredis.Client) *PresenceStore {
	return &PresenceStore{client: client}
}

func NewClient(addr, password string, db int) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func NewClientFromURL(redisURL string) (*goredis.Client, error) {
	options, err := goredis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return goredis.NewClient(options), nil
}

func (s *PresenceStore) Save(ctx context.Context, record presence.Record, ttl time.Duration) error {
	if ttl <= 0 {
		return s.Delete(ctx, record.UserID)
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, presenceKey(record.UserID), payload, ttl).Err()
}

func (s *PresenceStore) Get(ctx context.Context, userID string) (presence.Record, bool, error) {
	value, err := s.client.Get(ctx, presenceKey(userID)).Bytes()
	if err == goredis.Nil {
		return presence.Record{}, false, nil
	}
	if err != nil {
		return presence.Record{}, false, err
	}

	var record presence.Record
	if err := json.Unmarshal(value, &record); err != nil {
		return presence.Record{}, false, fmt.Errorf("decode presence %q: %w", presenceKey(userID), err)
	}
	return record, true, nil
}

func (s *PresenceStore) Delete(ctx context.Context, userID string) error {
	return s.client.Del(ctx, presenceKey(userID)).Err()
}

func (s *PresenceStore) List(ctx context.Context) ([]presence.Record, error) {
	var cursor uint64
	var records []presence.Record

	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, presenceKeyPrefix+"*", 100).Result()
		if err != nil {
			return nil, err
		}
		cursor = nextCursor

		for _, key := range keys {
			value, err := s.client.Get(ctx, key).Bytes()
			if err == goredis.Nil {
				continue
			}
			if err != nil {
				return nil, err
			}

			var record presence.Record
			if err := json.Unmarshal(value, &record); err != nil {
				return nil, fmt.Errorf("decode presence %q: %w", key, err)
			}
			records = append(records, record)
		}

		if cursor == 0 {
			break
		}
	}

	return records, nil
}

func presenceKey(userID string) string {
	return presenceKeyPrefix + strings.NewReplacer(":", "_", "/", "_", "\\", "_").Replace(userID)
}
