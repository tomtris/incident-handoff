package main

import (
	"context"
	"strconv"
	"sync"
	"time"
)

type InMemoryOnCallStore struct {
	mu            sync.RWMutex
	OnCallEntries map[string]OnCallShiftEntry
	currentID     int
}

func NewInMemoryOnCallStore() (*InMemoryOnCallStore, error) {
	s := InMemoryOnCallStore{
		OnCallEntries: make(map[string]OnCallShiftEntry), // id:Entry
		currentID:     0,
	}
	return &s, nil
}

func (s *InMemoryOnCallStore) Create(ctx context.Context, entry OnCallShiftEntry) (OnCallShiftEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentID++
	entry.ID = OnCallShiftEntryIDPrefix + strconv.Itoa(s.currentID)
	s.OnCallEntries[entry.ID] = entry
	return entry, nil
}

func (s *InMemoryOnCallStore) CurrentOnCall(ctx context.Context, service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	for _, each := range s.OnCallEntries {
		// startsat <= now  AND  now > EndsAt
		if each.Service == service && !each.StartsAt.After(now) && each.EndsAt.After(now) {
			return each.Username, nil
		}
	}
	return "", OnCallShiftEntryNotFound
}

func (s *InMemoryOnCallStore) ListOnCalls(ctx context.Context, startsAt *time.Time, endsAt *time.Time) ([]OnCallShiftEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := []OnCallShiftEntry{}
	for _, each := range s.OnCallEntries {
		// startsat >= startsAt
		if startsAt != nil && each.StartsAt.Before(*startsAt) {
			continue
		}
		// startsat < endsAt
		if endsAt != nil && !each.StartsAt.Before(*endsAt) {
			continue
		}
		entries = append(entries, each)
	}
	return entries, nil
}

func (s *InMemoryOnCallStore) UpdateOnCall(ctx context.Context, updatedEntry OnCallShiftEntry) (OnCallShiftEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.OnCallEntries[updatedEntry.ID]; !ok {
		return OnCallShiftEntry{}, OnCallShiftEntryNotFound
	}

	s.OnCallEntries[updatedEntry.ID] = updatedEntry
	return updatedEntry, nil
}
