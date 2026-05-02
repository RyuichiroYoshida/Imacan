package presence

import (
	"context"
	"errors"
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
	store             Store
	classTTL          time.Duration
	schoolCloseHour   int
	schoolCloseMinute int
}

func NewService(store Store, classTTL time.Duration, schoolCloseHour, schoolCloseMinute int) *Service {
	return &Service{
		store:             store,
		classTTL:          classTTL,
		schoolCloseHour:   schoolCloseHour,
		schoolCloseMinute: schoolCloseMinute,
	}
}

func (s *Service) Update(ctx context.Context, userID string, activity generated.Activity, now time.Time) (Record, bool, error) {
	if !activity.Valid() {
		return Record{}, false, ErrInvalidActivity
	}

	if activity == generated.OUT {
		if err := s.store.Delete(ctx, userID); err != nil {
			return Record{}, false, err
		}
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
	if err := s.store.Save(ctx, record, expiresAt.Sub(now)); err != nil {
		return Record{}, false, err
	}

	return record, true, nil
}

func (s *Service) Summary(ctx context.Context, now time.Time) (Summary, error) {
	records, err := s.store.List(ctx)
	if err != nil {
		return Summary{}, err
	}

	var summary Summary
	for _, record := range records {
		if !record.ExpiresAt.After(now) {
			_ = s.store.Delete(ctx, record.UserID)
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

	return summary, nil
}

func (s *Service) Current(ctx context.Context, userID string, now time.Time) (Record, bool, error) {
	record, ok, err := s.store.Get(ctx, userID)
	if err != nil || !ok {
		return Record{}, false, err
	}

	if !record.ExpiresAt.After(now) {
		if err := s.store.Delete(ctx, userID); err != nil {
			return Record{}, false, err
		}
		return Record{}, false, nil
	}

	return record, true, nil
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
