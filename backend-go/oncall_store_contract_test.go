package main

import (
	"context"
	"testing"
	"time"
)

// onCallStoreContract exercises the OnCallStore interface. Any implementation
// must pass it. newStore must return a fresh, empty store on each call —
// including a reset ID counter, so the ID-sequence assertion holds.
func onCallStoreContract(t *testing.T, newStore func() OnCallStore) {
	ctx := context.Background()
	base := time.Now().UTC().Truncate(time.Millisecond)

	t.Run("Create returns sequential IDs and preserves fields", func(t *testing.T) {
		store := newStore()
		req := OnCallShiftEntry{
			Service:  "payment",
			Username: "tom",
			StartsAt: base.Add(-1 * time.Minute),
			EndsAt:   base.Add(100 * time.Minute),
		}
		e1, err := store.Create(ctx, req)
		if err != nil {
			t.Fatalf("create 1: %v", err)
		}
		e2, err := store.Create(ctx, req)
		if err != nil {
			t.Fatalf("create 2: %v", err)
		}
		if e1.ID != OnCallShiftEntryIDPrefix+"1" || e2.ID != OnCallShiftEntryIDPrefix+"2" {
			t.Fatalf("expected sequential IDs, got %q %q", e1.ID, e2.ID)
		}
		if e1.Service != req.Service || e1.Username != req.Username {
			t.Fatal("service/username not preserved")
		}
		if !e1.StartsAt.Equal(req.StartsAt) || !e1.EndsAt.Equal(req.EndsAt) {
			t.Fatal("timestamps not preserved")
		}
	})

	t.Run("CurrentOnCall finds active shift", func(t *testing.T) {
		store := newStore()
		_, err := store.Create(ctx, OnCallShiftEntry{
			Service:  "payment",
			Username: "tom",
			StartsAt: base.Add(-1 * time.Minute),
			EndsAt:   base.Add(100 * time.Minute),
		})
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
		got, err := store.CurrentOnCall(ctx, "payment")
		if err != nil {
			t.Fatalf("current: %v", err)
		}
		if got != "tom" {
			t.Fatalf("expected tom, got %q", got)
		}
	})

	t.Run("CurrentOnCall not found for unknown service", func(t *testing.T) {
		store := newStore()
		_, err := store.CurrentOnCall(ctx, "not-exist")
		if err != OnCallShiftEntryNotFound {
			t.Fatalf("expected OnCallShiftEntryNotFound, got %v", err)
		}
	})

	t.Run("CurrentOnCall excludes expired and future shifts", func(t *testing.T) {
		store := newStore()
		if _, err := store.Create(ctx, OnCallShiftEntry{
			Service: "web", Username: "old",
			StartsAt: base.Add(-2 * time.Hour), EndsAt: base.Add(-1 * time.Hour),
		}); err != nil {
			t.Fatalf("seed expired: %v", err)
		}
		if _, err := store.Create(ctx, OnCallShiftEntry{
			Service: "web", Username: "next",
			StartsAt: base.Add(1 * time.Hour), EndsAt: base.Add(2 * time.Hour),
		}); err != nil {
			t.Fatalf("seed future: %v", err)
		}
		_, err := store.CurrentOnCall(ctx, "web")
		if err != OnCallShiftEntryNotFound {
			t.Fatalf("expected not found, got %v", err)
		}
	})

	t.Run("ListOnCalls no bounds returns all", func(t *testing.T) {
		store := newStore()
		seedThree(t, ctx, store, base)
		got, err := store.ListOnCalls(ctx, nil, nil)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(got) != 3 {
			t.Fatalf("expected 3, got %d", len(got))
		}
	})

	t.Run("ListOnCalls lower bound inclusive", func(t *testing.T) {
		store := newStore()
		seedThree(t, ctx, store, base) // StartsAt at base+0h, +1h, +2h
		from := base.Add(1 * time.Hour)
		got, err := store.ListOnCalls(ctx, &from, nil)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2, got %d", len(got))
		}
	})

	t.Run("ListOnCalls upper bound exclusive", func(t *testing.T) {
		store := newStore()
		seedThree(t, ctx, store, base)
		to := base.Add(2 * time.Hour)
		got, err := store.ListOnCalls(ctx, nil, &to)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2, got %d", len(got))
		}
	})

	t.Run("ListOnCalls both bounds half-open window", func(t *testing.T) {
		store := newStore()
		seedThree(t, ctx, store, base)
		from := base.Add(1 * time.Hour)
		to := base.Add(2 * time.Hour)
		got, err := store.ListOnCalls(ctx, &from, &to)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1, got %d", len(got))
		}
		if got[0].StartsAt.Equal(base.Add(1*time.Hour)) == false {
			t.Fatalf("wrong entry: StartsAt %v", got[0].StartsAt)
		}
	})

	t.Run("UpdateOnCall replaces existing", func(t *testing.T) {
		store := newStore()
		created, err := store.Create(ctx, OnCallShiftEntry{
			Service: "payment", Username: "tom",
			StartsAt: base, EndsAt: base.Add(1 * time.Hour),
		})
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
		created.Username = "jerry"
		updated, err := store.UpdateOnCall(ctx, created)
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		if updated.Username != "jerry" {
			t.Fatalf("expected jerry, got %q", updated.Username)
		}
		got, err := store.ListOnCalls(ctx, nil, nil)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(got) != 1 || got[0].Username != "jerry" {
			t.Fatalf("update not persisted: %+v", got)
		}
	})

	t.Run("UpdateOnCall not found for unknown ID", func(t *testing.T) {
		store := newStore()
		_, err := store.UpdateOnCall(ctx, OnCallShiftEntry{
			ID: OnCallShiftEntryIDPrefix + "does-not-exist", Username: "ghost",
		})
		if err != OnCallShiftEntryNotFound {
			t.Fatalf("expected OnCallShiftEntryNotFound, got %v", err)
		}
	})
}

func seedThree(t *testing.T, ctx context.Context, store OnCallStore, base time.Time) {
	t.Helper()
	for i, off := range []time.Duration{0, time.Hour, 2 * time.Hour} {
		_, err := store.Create(ctx, OnCallShiftEntry{
			Service:  "web",
			Username: "u" + string(rune('a'+i)),
			StartsAt: base.Add(off),
			EndsAt:   base.Add(off + 30*time.Minute),
		})
		if err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
	}
}
