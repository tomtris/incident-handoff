package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestHealthCheck(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/healthz", nil)
	ctx := context.WithValue(req.Context(), requestIDKey, "test-health-check-id")
	healthCheck(rec, req.WithContext(ctx))
	if rec.Code != 200 {
		t.Fatalf("Expected Code %v, get %v", 200, rec.Code)
	}
	var body map[string]string
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("expected no error, got err %v", err.Error())
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status %v, got %v", "ok", body["status"])
	}
}

func TestReadyCheckInMemory(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/readyz", nil)
	ctx := context.WithValue(req.Context(), requestIDKey, "test-ready-check-id")
	readyCheck(nil)(rec, req.WithContext(ctx))
	if rec.Code != 200 {
		t.Fatalf("Expected Code %v, get %v", 200, rec.Code)
	}
	var body map[string]string
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("expected no error, got err %v", err.Error())
	}
	if body["status"] != "ready" {
		t.Fatalf("expected status %v, got %v", "ready", body["status"])
	}
}

func TestReadyCheckDBUnavailable(t *testing.T) {
	clientOpts := options.Client().
		ApplyURI("mongodb://127.0.0.1:1"). // 1-1023 are privileged -> closed automatically
		SetServerSelectionTimeout(500 * time.Millisecond)
	client, err := mongo.Connect(clientOpts) // Just spawn monitoring goroutine directly and returns success. only fail on bad configuration (e.g URL)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer client.Disconnect(context.Background())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/readyz", nil)
	ctx := context.WithValue(req.Context(), requestIDKey, "test-ready-db-down")
	readyCheck(client.Database("test"))(rec, req.WithContext(ctx))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected %d, got %d", http.StatusServiceUnavailable, rec.Code)
	}
	var body map[string]ErrorMessageJSON
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	e := body["error"]
	if e.ErrorCode != "DATABASE_UNAVAILABLE" {
		t.Fatalf("expected DATABASE_UNAVAILABLE, got %q", e.ErrorCode)
	}
	if e.RequestID != "test-ready-db-down" {
		t.Fatalf("expected request id propagated, got %q", e.RequestID)
	}
}

func TestReadyCheckDBAvailable(t *testing.T) {
	clientOpts := options.Client().
		ApplyURI("mongodb://127.0.0.1:27017/?directConnection=true").
		SetServerSelectionTimeout(500 * time.Millisecond)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer client.Disconnect(context.Background())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/readyz", nil)
	ctx := context.WithValue(req.Context(), requestIDKey, "test-ready-db-down")
	readyCheck(client.Database("test"))(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}
}
