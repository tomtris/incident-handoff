package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	TRIGGERED     = "triggered"
	ACKNOWLEDGED  = "acknowledged"
	INVESTIGATING = "investigating"
	MITIGATED     = "mitigated"
	RESOLVED      = "resolved"
)

const requestIDKey = "request_id"

type Incident struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Service   string          `json:"service"`
	Severity  string          `json:"severity"` // SEV1, SEV2, SEV3
	Status    string          `json:"status"`   // triggered, acknowledged, investigating, mitigated, resolved
	OpenedBy  string          `json:"opened_by"`
	OnCall    string          `json:"on_call"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Entries   []TimelineEntry `json:"entries"`
}

type TimelineEntry struct {
	ID     string    `json:"id"`
	Time   time.Time `json:"time"`
	Author string    `json:"author"`
	Type   string    `json:"type"` // observation, action, discovery, open_question, state_change
	Text   string    `json:"text"`
}

var incidents []Incident
var nextIncidentID int
var nextEntryTimelineID int

// type IncidentOption func(*Incident)

// func withOncall(s string) func(*Incident) {
// 	return func(i *Incident) { i.OnCall = s }
// }

func defaultNewIncident() Incident {
	nextIncidentID++
	return Incident{
		ID:        "inc-" + strconv.Itoa(nextIncidentID),
		Status:    TRIGGERED,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createNewIncident(w http.ResponseWriter, r *http.Request) {
	newIncident := defaultNewIncident()
	err := json.NewDecoder(r.Body).Decode(&newIncident)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	incidents = append(incidents, newIncident)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newIncident)
	return
}

func defaultNewTimelineEntry() TimelineEntry {
	nextEntryTimelineID++
	return TimelineEntry{
		ID:   "ent-" + strconv.Itoa(nextEntryTimelineID),
		Time: time.Now(),
	}
}

func addTimelineEntry(w http.ResponseWriter, r *http.Request) {
	newTimelineEntry := defaultNewTimelineEntry()
	err := json.NewDecoder(r.Body).Decode(&newTimelineEntry)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	incidentId := r.PathValue("id")
	for idx, _ := range incidents {
		if incidents[idx].ID == incidentId {
			incidents[idx].Entries = append(incidents[idx].Entries, newTimelineEntry)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newTimelineEntry)
			return
		}
	}
	http.Error(w, "incident not found", http.StatusNotFound)
}

func listAllIncidents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "allication/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incidents)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /incidents", createNewIncident)
	mux.HandleFunc("POST /incidents/{id}/entries", addTimelineEntry)
	mux.HandleFunc("GET /incidents", listAllIncidents)
	handler := RequestIDMiddleware(LoggingMiddleware(CORSMiddleware(mux)))
	// mux.HandleFunc("GET /incidents", listOneIncident)
	// mux.HandleFunc("GET /incidents/{id}/handoff", listHandoff)
	// mux.HandleFunc("GET /healthz", healthCheck)
	// mux.HandleFunc("PATCH /incidents", UpdateOneIncident)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
