package handlers

import (
	"encoding/json"
	"net/http"
)

func (e Events) JSONError(w http.ResponseWriter, je interface{}, code int) {
    w.Header().Set("content-type", "application/json; charset=utf-8")
    w.Header().Set("x-content-type-options", "nosniff")
    w.WriteHeader(code)
    enc := json.NewEncoder(w)
    err := enc.Encode(je)

    if err != nil {
        e.l.Printf("Error marshalling JSON error: %v", err)
    }
}

func (a Auth) JSONError(w http.ResponseWriter, je interface{}, code int) {
    w.Header().Set("content-type", "application/json; charset=utf-8")
    w.Header().Set("x-content-type-options", "nosniff")
    w.WriteHeader(code)
    enc := json.NewEncoder(w)
    err := enc.Encode(je)

    if err != nil {
        a.l.Printf("Error marshalling JSON error: %v", err)
    }
}
