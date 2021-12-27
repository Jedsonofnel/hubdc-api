package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Jedsonofnel/hubdc-api/data"
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

    if r.Method == http.MethodPost {
        e.addEvent(rw, r)
        return
    }

    if r.Method == http.MethodPut {
        e.l.Println("PUT")
        // extract id from url path
        reg := regexp.MustCompile(`/([0-9]+)`)
        g := reg.FindAllStringSubmatch(r.URL.Path, -1)

        if len(g) != 1 {
            http.Error(rw, "Invalid URL", http.StatusBadRequest)
            return
        }

        if len(g[0]) != 2 {
            http.Error(rw, "Invalid URL", http.StatusBadRequest)
            return
        }
        idString := g[0][1]
        id, err := strconv.Atoi(idString)
        if err != nil {
            http.Error(rw, "Invalid URL", http.StatusBadRequest)
            return
        }

        e.updateEvents(id, rw, r)
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

func (e *Events) addEvent(rw http.ResponseWriter, r *http.Request) {
    e.l.Println("Handle POST Event")

    event := &data.Event{}

    err := event.FromJSON(r.Body)
    if err != nil {
        http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
    }

    data.AddEvent(event)
}

func (e *Events) updateEvents (id int, rw http.ResponseWriter, r *http.Request) {
    e.l.Println("Handle PUT Event")

    event := &data.Event{}

    err := event.FromJSON(r.Body)
    if err != nil {
        http.Error(rw, "unable to unmarhsal json", http.StatusBadRequest)
    }

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
