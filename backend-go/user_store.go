package main

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserStore interface {
	Create(ctx context.Context, u User) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
}

func NewUserStore(ctx context.Context, db *mongo.Database) (UserStore, error) {
	if db == nil {
		slog.Info("use in-memory store for UserStore")
		return NewMemoryUserStore()
	}
	slog.Info("use MongoStore for UserStore")
	return NewMongoUserStore(ctx, db.Collection(CollectionUsers))
}
