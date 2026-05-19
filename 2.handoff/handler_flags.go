package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (incHandler *IncidentHandler) CreateFlag(w http.ResponseWriter, r *http.Request) {
	f := FeatureFlag{}

	err := json.NewDecoder(r.Body).Decode(&f)
	requestID := r.Context().Value(requestIDKey).(string)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: BAD_REQUEST,
			Message:   ErrBadRequest.Error(),
			RequestID: requestID,
		})
		return
	}
	err = f.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: BAD_REQUEST,
			Message:   err.Error(),
			RequestID: requestID,
		})
		return
	}
	incHandler.FlagStore.Create(f)
	writeJSON(w, http.StatusCreated, requestID, f)
}

func (incHandler *IncidentHandler) UpdateFlag(w http.ResponseWriter, r *http.Request) {
	u := FeatureFlagUpdate{}

	requestID := r.Context().Value(requestIDKey).(string)
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: BAD_REQUEST,
			Message:   ErrBadRequest.Error(),
			RequestID: requestID,
		})
		return
	}

	u.Name = r.PathValue("name")
	err = u.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: BAD_REQUEST,
			Message:   err.Error(),
			RequestID: requestID,
		})
		return
	}

	// flag Store
	err = incHandler.FlagStore.Update(u)
	if errors.Is(err, ErrFlagNotfound) {
		writeError(w, http.StatusNotFound, ErrorMessageJSON{
			ErrorCode: FLAG_NOT_FOUND,
			Message:   ErrFlagNotfound.Error(),
			RequestID: requestID,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (incHandler *IncidentHandler) ListAllFlag(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(requestIDKey).(string)
	allFlags := incHandler.FlagStore.AllFlags()
	writeJSON(w, http.StatusOK, requestID, allFlags)
}

func (incHandler *IncidentHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(requestIDKey).(string)
	flagName := r.PathValue("name")
	userID := r.URL.Query().Get("user_id")

	// Validate
	if userID == "" {
		writeError(w, http.StatusNotFound, ErrorMessageJSON{
			ErrorCode: BAD_REQUEST,
			Message:   ErrBadRequest.Error(),
			RequestID: requestID,
		})
		return
	}

	// evaluate and answer
	flagAnswer, err := incHandler.FlagStore.Evaluate(flagName, userID)
	if err != nil {
		if errors.Is(err, ErrFlagNotfound) {
			writeError(w, http.StatusNotFound, ErrorMessageJSON{
				ErrorCode: FLAG_NOT_FOUND,
				Message:   ErrBadRequest.Error(),
				RequestID: requestID,
			})
		} else {
			writeError(w, http.StatusInternalServerError, ErrorMessageJSON{
				ErrorCode: INTERNAL_SERVER_ERROR,
				Message:   err.Error(),
				RequestID: requestID,
			})
		}
		return
	}

	writeJSON(w, http.StatusOK, requestID, flagAnswer)
}
