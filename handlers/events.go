package handlers

import (
	"context"
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
    e.l.Println("Handling INDEX request")

	le, err := e.s.GetEvents()
    if err != nil {
        e.l.Printf("Error handling INDEX request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error accessing database"),
            http.StatusInternalServerError,
        )
        return
    }

    rw.Header().Add("content-type", "application/json; charset=utf-8")
    rw.WriteHeader(http.StatusOK)
	err = le.ToJSON(rw)
	if err != nil {
        e.l.Printf("Error handling INDEX request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error marshaling data to JSON for response"),
            http.StatusInternalServerError,
        )
        return
	}
}

func (e Events) Upcoming(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("Handling UPCOMING request")

    le, err := e.s.GetUpcomingEvents()
    if err != nil   {
        e.l.Printf("Error handling UPCOMING request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error accessing database"),
            http.StatusInternalServerError,
        )
        return
    }

    rw.Header().Add("content-type", "application/json; charset=utf-8")
    rw.WriteHeader(http.StatusOK)
    err = le.ToJSON(rw)
    if err != nil {
        e.l.Printf("Error handling UPCOMING request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error marshaling data to JSON for response"),
            http.StatusInternalServerError,
        )
        return
    }
}

func(e Events) Show(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("Handling SHOW request")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        e.l.Printf("Error handling SHOW request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Unable to convert id to integer"),
            http.StatusBadRequest,
        )
        return
    }

    le, err := e.s.GetEvent(id)
    if err != nil{
        e.l.Printf("Error handling SHOW request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Event not found"),
            http.StatusNotFound,
        )
        return
    }

    rw.Header().Add("content-type", "application/json; charset=utf-8")
    rw.WriteHeader(http.StatusOK)
    err = le.ToJSON(rw)
	if err != nil {
        e.l.Printf("Error handling SHOW request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error marshaling data to JSON for response"),
            http.StatusInternalServerError,
        )
        return
	}
}

func (e *Events) Create(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("Handling CREATE request")

    event :=r.Context().Value(KeyEvent{}).(*data.Event)
    ret, err := e.s.CreateEvent(event)
    if err != nil {
        e.l.Printf("Error handling CREATE request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error creating resource"),
            http.StatusInternalServerError,
        )
        return
    }

    // returning created event
    rw.Header().Add("content-type", "application/json; charset=utf-8")
    rw.WriteHeader(http.StatusOK)
    err = ret.ToJSON(rw)
	if err != nil {
        e.l.Printf("Error handling CREATE request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error marshaling data to JSON for response"),
            http.StatusInternalServerError,
        )
        return
	}
}

func (e Events) Update(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("Handling UPDATE request")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        e.l.Printf("Error handling UPDATE request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Unable to convert id to int"),
            http.StatusBadRequest,
        )
        return
    }

    event := r.Context().Value(KeyEvent{}).(*data.Event)
    ret, err := e.s.UpdateEvent(id, event)
    if err != nil {
        e.l.Printf("Error handling UPDATE request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error updating event"),
            http.StatusBadRequest,
        )
        return
    }

    // return updated event to client
    rw.Header().Add("content-type", "application/json; charset=utf-8")
    err = ret.ToJSON(rw)
	if err != nil {
        e.l.Printf("Error handling UPDATE request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Error marshaling data to JSON for response"),
            http.StatusInternalServerError,
        )
        return
	}
}

func (e Events) Delete(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("Handling DELETE request")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        e.l.Printf("Error handling DELETE request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Unable to convert id to int"),
            http.StatusBadRequest,
        )
        return
    }

    err = e.s.DeleteEvent(id)
    if err != nil {
        e.l.Printf("Error handling DELETE request: %v", err)
        e.JSONError(
            rw,
            data.NewJEs("Event not found"),
            http.StatusNotFound,
        )
        return
    }
}

type KeyEvent struct {}

func (e Events) MiddlewareEventValidation(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        event := &data.Event{}

        err := event.FromJSON(r.Body)
        if err != nil {
            e.JSONError(
                rw,
                data.NewJEs("Error marshaling data to JSON for response"),
                http.StatusInternalServerError,
            )
            return
        }

        // validate json
        errs := event.Validate()
        if errs.Errors != nil {
            e.JSONError(
                rw,
                errs,
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
