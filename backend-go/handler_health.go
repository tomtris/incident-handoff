package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, getRequestID(r), map[string]string{"status": "ok"})
}

func readyCheck(mongoDB *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// No DB -> use Memory -> always ready
		if mongoDB == nil {
			writeJSON(w, http.StatusOK, getRequestID(r), map[string]string{"status": "ready"})
			return
		}
		// Check health with mongoDB
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err := mongoDB.RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err()
		if err != nil {
			requestID := getRequestID(r)
			slog.Error("DB not available", "requestID", requestID, "error", err)
			writeError(w, http.StatusServiceUnavailable, ErrorMessageJSON{
				ErrorCode: "DATABASE_UNAVAILABLE",
				Message:   "database unavailable",
				RequestID: requestID,
			})
			return
		}
		writeJSON(w, http.StatusOK, getRequestID(r), map[string]string{"status": "ready"})
	}
}
