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
    s data.EventStore
}

func NewEvents(l *log.Logger, s data.EventStore) *Events {
    return &Events{l, s}
}

func (e Events) Index(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("handle INDEX request")

	le, err := e.s.GetEvents()
    if err != nil {
        http.Error(rw, "Error accessing database", http.StatusInternalServerError)
    }

    rw.Header().Add("content-type", "application/json; charset=utf-8")
	err = le.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (e Events) Upcoming(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("handle UPCOMING request")

    le, err := e.s.GetUpcomingEvents()
    if err != nil   {
        e.l.Println(err)
        http.Error(rw, "Error accessing database", http.StatusInternalServerError)
    }

    rw.Header().Add("content-type", "application/json; charset=utf-8")
    err = le.ToJSON(rw)
    if err != nil {
        http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
    }
}

func(e Events) Show(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("handle SHOW request")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(rw, "Unable to convert id", http.StatusBadRequest)
        return
    }

    le, err := e.s.GetEvent(id)
    if err != nil{
        http.Error(rw, "Event not found", http.StatusNotFound)
        return
    }

    rw.Header().Add("content-type", "application/json; charset=utf-8")
    err = le.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
        return
	}
}

func (e *Events) Create(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("handle CREATE request")

    event :=r.Context().Value(KeyEvent{}).(*data.Event)
    ret, err := e.s.CreateEvent(event)
    if err != nil {
        http.Error(rw, "Error creating resource: %v", http.StatusInternalServerError)
    }

    // returning created event
    rw.Header().Add("content-type", "application/json; charset=utf-8")
    err = ret.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
        return
	}
}

func (e Events) Update(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("handle UPDATE request")
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(rw, "Unable to convert id", http.StatusBadRequest)
        return
    }

    event := r.Context().Value(KeyEvent{}).(*data.Event)

    ret, err := e.s.UpdateEvent(id, event)
    if err != nil {
        http.Error(rw, "Error updating events", http.StatusInternalServerError)
        return
    }

    // return updated event to client
    rw.Header().Add("content-type", "application/json; charset=utf-8")
    err = ret.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
        return
	}
}

func (e Events) Delete(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("handle DELETE request")
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(rw, "Unable to convert id", http.StatusBadRequest)
        return
    }

    err = e.s.DeleteEvent(id)
    if err != nil {
        http.Error(rw, "Event not found", http.StatusNotFound)
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
