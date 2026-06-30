package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoOnCallStore struct {
	col *mongo.Collection
}

func NewMongoOnCallStore(ctx context.Context, col *mongo.Collection) (*MongoOnCallStore, error) {
	return &MongoOnCallStore{col: col}, nil
}

func (s *MongoOnCallStore) Create(ctx context.Context, entry OnCallShiftEntry) (OnCallShiftEntry, error) {
	id, err := mongoNextID(ctx, CollectionOnCallShifts, OnCallShiftEntryIDPrefix)
	if err != nil {
		return OnCallShiftEntry{}, fmt.Errorf("can not get next on-call id: %w", err)
	}
	entry.ID = id

	_, err = s.col.InsertOne(ctx, entry)
	if err != nil {
		return OnCallShiftEntry{}, fmt.Errorf("insert on-call error: %w", err)
	}
	return entry, nil
}

func (s *MongoOnCallStore) CurrentOnCall(ctx context.Context, service string) (string, error) {
	now := time.Now()

	filter := bson.M{
		"service":   service,
		"starts_at": bson.M{"$lte": now},
		"ends_at":   bson.M{"$gt": now},
	}

	var entry OnCallShiftEntry
	err := s.col.FindOne(ctx, filter).Decode(&entry)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", OnCallShiftEntryNotFound
		}
		return "", fmt.Errorf("current on-call query error: %w", err)
	}
	return entry.Username, nil
}

func (s *MongoOnCallStore) ListOnCalls(ctx context.Context, startsAt *time.Time, endsAt *time.Time) ([]OnCallShiftEntry, error) {
	startCond := bson.M{}
	if startsAt != nil {
		startCond["$gte"] = *startsAt
	}
	if endsAt != nil {
		startCond["$lt"] = *endsAt
	}

	filter := bson.M{}
	if len(startCond) > 0 {
		filter["starts_at"] = startCond
	}

	cursor, err := s.col.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list on-calls query error: %w", err)
	}

	entries := []OnCallShiftEntry{}
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, fmt.Errorf("list on-calls decode error: %w", err)
	}
	return entries, nil
}

func (s *MongoOnCallStore) UpdateOnCall(ctx context.Context, updatedEntry OnCallShiftEntry) (OnCallShiftEntry, error) {
	filter := bson.M{"_id": updatedEntry.ID}

	res, err := s.col.ReplaceOne(ctx, filter, updatedEntry)
	if err != nil {
		return OnCallShiftEntry{}, fmt.Errorf("update on-call error: %w", err)
	}
	if res.MatchedCount == 0 {
		return OnCallShiftEntry{}, OnCallShiftEntryNotFound
	}
	return updatedEntry, nil
}
