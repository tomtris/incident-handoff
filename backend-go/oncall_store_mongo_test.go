package main

import (
	"log"
	"testing"
)

func setupOnCallStoreMongoContractEnv(t *testing.T) *MongoOnCallStore {
	t.Helper()
	config := loadConfig()
	db = getMongoDatabase(config)
	db.Drop(t.Context())
	NewMongoOnCallStore, err := NewMongoOnCallStore(t.Context(), db.Collection(CollectionOnCallShifts))
	if err != nil {
		log.Fatalf("%v", err)

	}
	return NewMongoOnCallStore
}

func TestOnCallStoreMongoContract(t *testing.T) {
	onCallStoreContract(t, func() OnCallStore {
		return setupOnCallStoreMongoContractEnv(t)
	})
}
