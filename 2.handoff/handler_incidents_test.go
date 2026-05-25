package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshalNewEntryEvent(t *testing.T) {
	entryTimeline := TimelineEntry{
		ID:     "TLE-1",
		Time:   time.Now(),
		Author: "anh",
		Type:   OBSERVATION,
		Text:   "test entry",
	}
	rawMsg := marshalNewEntryEvent("INC-test1", entryTimeline)
	var event map[string]any
	json.Unmarshal(rawMsg, &event)

	if event["type"] != "new_entry" {
		t.Fatalf("type expected %v, get %v", OBSERVATION, event["type"])
	}
	if event["incident_id"] != "INC-test1" {
		t.Fatalf("incident_id %v, get %v", "INC-test1", event["incident_id"])
	}
	e := event["entry"].(map[string](any))
	if e["author"] != "anh" {
		t.Fatalf("author expected %v, get %v", "anh", e["author"])
	}
	if e["type"] != OBSERVATION {
		t.Fatalf("type expected %v, get %v", OBSERVATION, e["type"])
	}
}

func TestMarshalIncidentUpdateEvent(t *testing.T) {
	inc := Incident{
		ID:        "INC-test1",
		Title:     "test title",
		Service:   "test service",
		Severity:  "SEV1",
		Status:    TRIGGERED,
		OpenedBy:  "anh",
		OnCall:    "tom",
		CreatedAt: time.Now().Add(-15 * time.Minute),
		UpdatedAt: time.Now(),
		Entries:   []TimelineEntry{},
	}

	rawMsg := marshalIncidentUpdateEvent(inc)
	var event map[string]any
	json.Unmarshal(rawMsg, &event)
	if event["type"] != "incident_updated" {
		t.Fatalf("type expected %v, get %v", "incident_updated", event["type"])
	}
	if event["type"] != "incident_updated" {
		t.Fatalf("type expected %v, get %v", "incident_updated", event["type"])
	}
	e := event["incident"].(map[string]any)
	if e["id"] != "INC-test1" {
		t.Fatalf("id expected %v, get %v", "INC-test1", e["id"])
	}
	if e["service"] != "test service" {
		t.Fatalf("id expected %v, get %v", "test service", e["service"])
	}
}
