package main

import "net/http"

func getRequestID(r *http.Request) string {
	return r.Context().Value(requestIDKey).(string)
}
