package main

import (
    "os"
    "encoding/json"
    "log"
    "io/ioutil"
    "sync"
    "time"
    "net/http"
)

type Event struct {
    Id      string  `json:"id"`
    When    string  `json:"when"`
    Where   string  `json:"where"`
    What    string  `json:"what"`
}

type eventHandler struct {
    sync.Mutex          `json:"-"`
    Store       []Event `json:"store"`
    StoreFile   string  `json:"-"`
    Password    string  `json:"-"`
}

func newEventHandler() *eventHandler {
    var e eventHandler
    e.StoreFile = "eventStore.json"

    // Open or create json "storeFile" if non-existent
    f, err := os.OpenFile(e.StoreFile, os.O_CREATE, 0644)
    if err != nil {
        log.Println(err)
    }
    defer f.Close()

    byteValue, err := ioutil.ReadAll(f)
    json.Unmarshal(byteValue, &e)

    // Set admin password needed to create new events
    password := os.Getenv("ADMIN_PASSWORD")
    if password == "" {
        panic("required env var ADMIN_PASSWORD not set")
    }
    e.Password = password

    return &e
}

func (e *Event) Time() time.Time {
    t, err := time.Parse("15:04 02-01-06", e.When)
    if err != nil {
        log.Println("error parsing time: ", err)
    }
    return t
}

func (h *eventHandler) SerialiseBaby() error {
    writeBytes, err := json.MarshalIndent(h, "", " ")
    if err != nil {
        return newHTTPError(
            err,
            "error writing db to file",
            http.StatusInternalServerError,
        )
    }
    os.WriteFile(h.StoreFile, writeBytes, 0644)
    return nil
}
// Return a slice of the next three events in order of soon -> latest
func (h *eventHandler) Upcoming() []Event {
    var next Event
    var nextNext Event
    var nextNextNext Event

    // Used to see if event time is today or afterwards
    lastMidnight := time.Now().Truncate(24*time.Hour)

    // Super janky way of finding next three events
    sortable := h.Store
    for len(sortable) != 0 {
        nextTimeSorts := []Event{}
        for _, e := range sortable {
            if !e.Time().After(lastMidnight) {
                continue
            }
            switch {
            case next.When == "" || e.Time().Before(next.Time()):
                if next.When != "" {
                    nextTimeSorts = append(nextTimeSorts, next)
                }
                next = e
            case nextNext.When == "" || e.Time().Before(nextNext.Time()):
                if nextNext.When != "" {
                    nextTimeSorts = append(nextTimeSorts, nextNext)
                }
                nextNext = e
            case nextNextNext.When == "" || e.Time().Before(nextNextNext.Time()):
                if nextNextNext.When != "" {
                    nextTimeSorts = append(nextTimeSorts, nextNextNext)
                }
                nextNextNext = e
            }
        }
        sortable = nextTimeSorts
    }

    // Find non empty events
    upcoming := []Event{}
    for _, e := range []Event{next, nextNext, nextNextNext} {
        if e.When != "" {
            upcoming = append(upcoming, e)
        }
    }
    return upcoming
}

func (h *eventHandler) FindWithID(id string) (Event, bool) {
    for _, e := range h.Store {
        if id == e.Id {
            return e, true
        }
    }
    return Event{}, false
}

func validEvent(reqEvent Event) (error, bool) {
    // Making sure all fields are entered
    missingFields := make([]string, 0)
    if reqEvent.When == "" {
        missingFields = append(missingFields, "when")
    }
    if reqEvent.Where == "" {
        missingFields = append(missingFields, "where")
    }
    if reqEvent.What == "" {
        missingFields = append(missingFields, "what")
    }
    if len(missingFields) != 0 {
        errString := "missing json fields: "
        for _, v := range missingFields {
            errString += v + " "
        }
        return newHTTPError(nil, errString, http.StatusBadRequest), false
    }

    // Test whether "when" is in the correct format
    _, err := time.Parse("15:04 02-01-06", reqEvent.When)
    if err != nil {
        return newHTTPError(
            err,
            "time not in '15:04 02-01-06' format",
            http.StatusBadRequest,
        ),
        false
    }
    return nil, true
}
