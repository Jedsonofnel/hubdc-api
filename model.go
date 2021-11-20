package main

import (
    "os"
    "encoding/json"
    "log"
    "io/ioutil"
    "sync"
    "time"
)

type Event struct {
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
