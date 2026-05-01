package presence

import (
	"context"
	"sync"
	"time"
)

type Store interface {
	Save(ctx context.Context, record Record, ttl time.Duration) error
	Delete(ctx context.Context, userID string) error
	List(ctx context.Context) ([]Record, error)
}

type MemoryStore struct {
	mu      sync.Mutex
	records map[string]Record
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		records: make(map[string]Record),
	}
}

func (s *MemoryStore) Save(ctx context.Context, record Record, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[record.UserID] = record
	return nil
}

func (s *MemoryStore) Delete(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.records, userID)
	return nil
}

func (s *MemoryStore) List(ctx context.Context) ([]Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	records := make([]Record, 0, len(s.records))
	for _, record := range s.records {
		records = append(records, record)
	}
	return records, nil
}
