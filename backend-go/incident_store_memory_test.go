package main

import "testing"

func TestStoreMemory(t *testing.T) {
	runStoreContractsTests(t, func(t *testing.T) IncidentStore {
		NewMemoryIncidentStore, _ := NewMemoryIncidentStore()
		return NewMemoryIncidentStore
	})
}
