package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jedsonofnel/hubdc-api/data"
)

func (e Events) JSONError(w http.ResponseWriter, jes data.JSONErrors, code int) {
    w.Header().Set("content-type", "application/json; charset=utf-8")
    w.Header().Set("x-content-type-options", "nosniff")
    w.WriteHeader(code)
    enc := json.NewEncoder(w)
    err := enc.Encode(jes)

    if err != nil {
        e.l.Printf("Error marshalling JSON error: %v", err)
    }
}

func (a Auth) JSONError(w http.ResponseWriter, jes data.JSONErrors, code int) {
    w.Header().Set("content-type", "application/json; charset=utf-8")
    w.Header().Set("x-content-type-options", "nosniff")
    w.WriteHeader(code)
    enc := json.NewEncoder(w)
    err := enc.Encode(jes)

    if err != nil {
        a.l.Printf("Error marshalling JSON error: %v", err)
    }
}
