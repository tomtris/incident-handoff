package main

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoUserStore struct {
	col *mongo.Collection
}

func NewMongoUserStore(ctx context.Context, col *mongo.Collection) (*MongoUserStore, error) {
	// Must create this index with unique true. So when creating user, no need to think about race.
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, fmt.Errorf("create username index: %v", err)
	}
	return &MongoUserStore{col: col}, nil
}

func (s *MongoUserStore) Create(ctx context.Context, u User) (User, error) {
	id, err := mongoNextID(ctx, CollectionCountersUser, UserIDPrefix)
	if err != nil {
		return User{}, fmt.Errorf("next user id: %w", err)
	}
	u.ID = id

	_, err = s.col.InsertOne(ctx, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return User{}, ErrUserAlreadyExist
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}
	return u, nil
}

func (s *MongoUserStore) GetByUsername(ctx context.Context, username string) (User, error) {
	var u User
	err := s.col.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}
