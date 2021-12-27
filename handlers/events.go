package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Jedsonofnel/hubdc-api/data"
	"github.com/gorilla/mux"
)

type Events struct {
	l *log.Logger
}

func NewEvents(l *log.Logger) *Events {
    return &Events{l}
}

func (e Events) GetEvents(rw http.ResponseWriter, r *http.Request) {
	le := data.GetEvents()
	err := le.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func(e Events) GetEvent(rw http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(rw, "Unable to convert id", http.StatusBadRequest)
        return
    }

    le, err := data.GetEvent(id)
    if err != nil{
        http.Error(rw, "Event not found", http.StatusNotFound)
        return
    }

    err = le.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
        return
	}
}

func (e *Events) AddEvent(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("Handle POST Event")

    event :=r.Context().Value(KeyEvent{}).(*data.Event)
    data.AddEvent(event)
}

func (e Events) UpdateEvent (rw http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(rw, "Unable to convert id", http.StatusBadRequest)
        return
    }

    e.l.Println("Handle PUT Event", id)
    event :=r.Context().Value(KeyEvent{}).(*data.Event)

    err = data.UpdateEvent(id, event)
    if err == data.ErrEventNotFound {
        http.Error(rw, "Event not found", http.StatusNotFound)
        return
    }

    if err != nil {
        http.Error(rw, "Event not found", http.StatusInternalServerError)
        return
    }
}

type KeyEvent struct {}

func (e Events) MiddlewareEventValidation(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        event := &data.Event{}

        err := event.FromJSON(r.Body)
        if err != nil {
            http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
            return
        }

        // validate json
        err = event.Validate()
        if err != nil {
            http.Error(
                rw,
                fmt.Sprintf("Error validating product %s", err),
                http.StatusBadRequest,
            )
            return
        }

        // use context to send worked data to next handler
        ctx := context.WithValue(r.Context(), KeyEvent{}, event)
        req := r.WithContext(ctx)

        next.ServeHTTP(rw, req)
    })
}
