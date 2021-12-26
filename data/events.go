package data

import (
    "encoding/json"
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

func GetEvents() Events {
    return eventList
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
