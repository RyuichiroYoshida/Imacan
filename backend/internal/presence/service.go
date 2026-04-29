package presence

import (
	"errors"
	"sync"
	"time"

	"github.com/RyuichiroYoshida/imacan/backend/internal/generated"
)

var ErrInvalidActivity = errors.New("invalid activity")

type Record struct {
	UserID    string
	Activity  generated.Activity
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type Summary struct {
	Total     int32
	Class     int32
	SelfStudy int32
}

type Service struct {
	mu                sync.Mutex
	records           map[string]Record
	classTTL          time.Duration
	schoolCloseHour   int
	schoolCloseMinute int
}

func NewService(classTTL time.Duration, schoolCloseHour, schoolCloseMinute int) *Service {
	return &Service{
		records:           make(map[string]Record),
		classTTL:          classTTL,
		schoolCloseHour:   schoolCloseHour,
		schoolCloseMinute: schoolCloseMinute,
	}
}

func (s *Service) Update(userID string, activity generated.Activity, now time.Time) (Record, bool, error) {
	if !activity.Valid() {
		return Record{}, false, ErrInvalidActivity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if activity == generated.OUT {
		delete(s.records, userID)
		return Record{
			UserID:    userID,
			Activity:  activity,
			UpdatedAt: now,
		}, false, nil
	}

	expiresAt := s.expiresAt(activity, now)
	record := Record{
		UserID:    userID,
		Activity:  activity,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
	}
	s.records[userID] = record

	return record, true, nil
}

func (s *Service) Summary(now time.Time) Summary {
	s.mu.Lock()
	defer s.mu.Unlock()

	var summary Summary
	for userID, record := range s.records {
		if !record.ExpiresAt.After(now) {
			delete(s.records, userID)
			continue
		}

		summary.Total++
		switch record.Activity {
		case generated.CLASS:
			summary.Class++
		case generated.SELFSTUDY:
			summary.SelfStudy++
		}
	}

	return summary
}

func (s *Service) expiresAt(activity generated.Activity, now time.Time) time.Time {
	closeAt := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		s.schoolCloseHour,
		s.schoolCloseMinute,
		0,
		0,
		now.Location(),
	)

	expiresAt := closeAt
	if activity == generated.CLASS {
		expiresAt = now.Add(s.classTTL)
		if expiresAt.After(closeAt) {
			expiresAt = closeAt
		}
	}

	if !expiresAt.After(now) {
		return now
	}
	return expiresAt
}
