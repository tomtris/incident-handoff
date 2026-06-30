package main

import (
	"log"
	"testing"
)

func TestOnCallStoreMemoryContract(t *testing.T) {
	onCallStoreContract(t, func() OnCallStore {
		InMemoryOnCallStore, err := NewInMemoryOnCallStore()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return InMemoryOnCallStore
	})
}
