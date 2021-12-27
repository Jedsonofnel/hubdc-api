package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type Event struct {
    ID      int     `json:"id"`
    What    string  `json:"what"`
    Where   string  `json:"where"`
    When    string  `json:"when"`
}

type Events []*Event

func (e *Events) ToJSON(w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(e)
}

func (e *Event) FromJSON(r io.Reader) error {
    d := json.NewDecoder(r)
    return d.Decode(e)
}

func GetEvents() Events {
    return eventList
}

func AddEvent(e *Event) {
    e.ID = getNextID()
    eventList = append(eventList, e)
}

func UpdateEvent(id int, e *Event) error {
    _, pos, err := findEvent(id)
    if err != nil {
        return err
    }

    e.ID = id
    eventList[pos] = e

    return nil
}

var ErrEventNotFound = fmt.Errorf("Event not found")

func findEvent(id int) (*Event, int, error) {
    for i, e := range eventList {
        if e.ID == id {
            return e, i, nil
        }
    }

    return nil, -1, ErrEventNotFound
}

func getNextID() int {
    le := eventList[len(eventList)-1]
    return le.ID + 1
}

var eventList = []*Event{
    {
        ID: 1,
        What: "Normal Hub Session",
        Where: "HPH",
        When: "17:15 15-01-22",
    },
    {
        ID: 2,
        What: "Lent Addresses",
        Where: "Chapel Close",
        When: "16:30 22-02-22",
    },
}
