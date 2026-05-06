package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IncidentHandler struct {
	Store Store
}

func (incHandler *IncidentHandler) CreateIncident(w http.ResponseWriter, r *http.Request) {
	newCreateIncidentRequest := CreateIncidentRequest{}
	err := json.NewDecoder(r.Body).Decode(&newCreateIncidentRequest)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: "BAD_REQUEST",
			Message:   fmt.Sprintf("Validation invalid: %s", err),
			RequestID: r.Context().Value(requestIDKey).(string),
		})
		return
	}

	err = newCreateIncidentRequest.Validate()
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: "BAD_REQUEST",
			Message:   fmt.Sprintf("Validation invalid: %s", err),
			RequestID: r.Context().Value(requestIDKey).(string),
		})
		return
	}

	createdIncident, err := incHandler.Store.CreateIncident(r.Context(), Incident{
		Title:    newCreateIncidentRequest.Title,
		Service:  newCreateIncidentRequest.Service,
		Severity: newCreateIncidentRequest.Severity,
		OpenedBy: newCreateIncidentRequest.OpenedBy,
	})

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorMessageJSON{
			ErrorCode: "INTERNAL_ERROR",
			Message:   "failed to create incident",
			RequestID: r.Context().Value(requestIDKey).(string),
		})
		return
	}

	writeJSON(w, http.StatusCreated, createdIncident)
	return
}
