package data

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// For quick use when pushing strings up to database expecting time.Time
type stringTime string

func (st stringTime) timeConv() time.Time {
    t, _ := time.Parse("15:04 02-01-06", string(st))
    return t
}

func (st *stringTime) parse(t time.Time) {
    *st = stringTime(t.Format("15:04 02-01-06"))
}

type Event struct {
    ID      int     `json:"id"`
    What    string  `json:"what" validate:"required"`
    Loc     string  `json:"loc" validate:"required"`
    When    stringTime  `json:"when" validate:"required,when"`
}

type Events []*Event

type EventStore struct {
    DB *sql.DB
}

func NewEventStore(url string) (EventStore, error) {
    var es EventStore
    db, err := sql.Open("postgres", url)
    if err != nil {
        return EventStore{}, err
    }
    es.DB = db

    err = es.DB.Ping()
    if err != nil {
        return EventStore{}, err
    }

    _, err = es.DB.Exec(createEventTable)
    if err != nil {
        return EventStore{}, err
    }

    return es, nil
}

func (es *EventStore) GetEvents() (Events, error) {
    rows, err := es.DB.Query(listEvents)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var events Events
    for rows.Next() {
        var e Event
        var et time.Time
        if err := rows.Scan(&e.ID, &e.What, &e.Loc, &et); err != nil {
            return nil, err
        }
        e.When.parse(et)
        events = append(events, &e)
    }
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
    return events, nil
}

func (es *EventStore) GetEvent(id int) (*Event, error) {
    row := es.DB.QueryRow(getEvent, id)
    var e Event
    var et time.Time
    err := row.Scan(&e.ID, &e.What, &e.Loc, &et)
    e.When.parse(et)
    return &e ,err
}

func (es *EventStore) CreateEvent(e *Event) (*Event, error) {
	row := es.DB.QueryRow(createEvent, e.What, e.Loc, e.When.timeConv())
	var re Event
    var et time.Time
	err := row.Scan(&re.ID, &re.What, &re.Loc, &et)
    re.When.parse(et)
	return &re, err
}

func (es *EventStore) UpdateEvent(id int, e *Event) (*Event, error) {
	row := es.DB.QueryRow(updateEvent, e.What, e.Loc, e.What, id)
	var re Event
    var et time.Time
	err := row.Scan(&e.ID, &re.What, &re.Loc, &et)
    re.When.parse(et)
	return &re, err
}

func (es *EventStore) DeleteEvent(id int) error {
    _, err := es.DB.Exec(deleteEvent, id)
    return err
}
