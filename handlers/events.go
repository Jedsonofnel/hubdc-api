package handlers

import (
	"github.com/Jedsonofnel/hubdc-api/data"
	"log"
	"net/http"
)

type Events struct {
	l *log.Logger
}

func NewEvents(l *log.Logger) *Events {
    return &Events{l}
}

func (e *Events) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        e.getEvents(rw, r)
        return
    }

    rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (e *Events) getEvents(rw http.ResponseWriter, r *http.Request) {
	le := data.GetEvents()
	err := le.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}
