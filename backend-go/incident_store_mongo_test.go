package main

import "testing"

const connectionDBString = "mongodb://127.0.0.1:27017/?directConnection=true"
const DBName = "incident_tracker"

func setupMongoTestEnv(t *testing.T) *MongoIncidentStore {
	t.Helper()
	config := loadConfig()
	db = getMongoDatabase(config)
	db.Drop(t.Context())
	mongoIncidentStore, _ := NewMongoIncidentStore(db.Collection(CollectionIncidents))
	return mongoIncidentStore
}

func TestMongoStore(t *testing.T) {
	runStoreContractsTests(t, func(t *testing.T) IncidentStore {
		return setupMongoTestEnv(t)
	})
}
