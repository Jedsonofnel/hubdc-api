package main

import (
    "os"
    // "encoding/json"
    "log"
    // "io/ioutil"
    // "sync"
    // "time"
    // "net/http"
    // "strconv"
    "database/sql"

    "github.com/joho/godotenv"
    _ "github.com/jackc/pgx/v4/stdlib"
)

type Event struct {
    Id      int  `json:"id"`
    When    string  `json:"when"`
    Where   string  `json:"where"`
    What    string  `json:"what"`
}

type eventHandler struct {
    Db          *sql.DB
    Password    string
    Username    string
}

func newEventHandler() *eventHandler {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %q", err)
    }
    log.Println(os.Getenv("DATABASE_URL"))

    var e eventHandler

    // Open connection to database
    db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }
    e.Db = db

    err = e.Db.Ping()
    if err != nil {
        log.Fatalf("Error pinging database: %q", err)
    } else {
        log.Println("Successfully connected to database")
    }

    query :=`CREATE TABLE IF NOT EXISTS event(
                event_id SERIAL PRIMARY KEY,
                event_when TIMESTAMP,
                event_where TEXT,
                event_what TEXT
            )`
    _, err = e.Db.Exec(query)
    if err != nil {
        log.Fatalf("Error creating table: %q", err)
    } else {
        log.Println("Table initialised")
    }

    // Set admin password needed to create new events
    e.Password = os.Getenv("ADMIN_PASSWORD")
    if e.Password == "" {
        log.Fatal("required env var ADMIN_PASSWORD not set!")
    }
    e.Username = os.Getenv("USERNAME")
    if e.Username == "" {
        log.Fatal("required env var USERNAME not set!")
    }

    return &e
}

// TODO Rewrite with new SQL system
// func (e *Event) Time() time.Time {
//     t, err := time.Parse("15:04 02-01-06", e.When)
//     if err != nil {
//         log.Println("error parsing time: ", err)
//     }
//     return t
// }
//
// func (h *eventHandler) GetBestID() string {
//     lowestId := 0
//     for _, e := range h.Store {
//         index, _ := strconv.Atoi(e.Id)
//         if index <= lowestId {
//             lowestId = index + 1
//         }
//     }
//     return strconv.Itoa(lowestId)
// }
// func (h *eventHandler) SerialiseBaby() error {
//     writeBytes, err := json.MarshalIndent(h, "", " ")
//     if err != nil {
//         return newHTTPError(
//             err,
//             "error writing db to file",
//             http.StatusInternalServerError,
//         )
//     }
//     os.WriteFile(h.StoreFile, writeBytes, 0644)
//     return nil
// }
//
// // Return a slice of the next three events in order of soon -> latest
// func (h *eventHandler) Upcoming() []Event {
//     var next Event
//     var nextNext Event
//     var nextNextNext Event
//
//     // Used to see if event time is today or afterwards
//     lastMidnight := time.Now().Truncate(24*time.Hour)
//
//     // Super janky way of finding next three events
//     sortable := h.Store
//     for len(sortable) != 0 {
//         nextTimeSorts := []Event{}
//         for _, e := range sortable {
//             if !e.Time().After(lastMidnight) {
//                 continue
//             }
//             switch {
//             case next.When == "" || e.Time().Before(next.Time()):
//                 if next.When != "" {
//                     nextTimeSorts = append(nextTimeSorts, next)
//                 }
//                 next = e
//             case nextNext.When == "" || e.Time().Before(nextNext.Time()):
//                 if nextNext.When != "" {
//                     nextTimeSorts = append(nextTimeSorts, nextNext)
//                 }
//                 nextNext = e
//             case nextNextNext.When == "" || e.Time().Before(nextNextNext.Time()):
//                 if nextNextNext.When != "" {
//                     nextTimeSorts = append(nextTimeSorts, nextNextNext)
//                 }
//                 nextNextNext = e
//             }
//         }
//         sortable = nextTimeSorts
//     }
//
//     // Find non empty events
//     upcoming := []Event{}
//     for _, e := range []Event{next, nextNext, nextNextNext} {
//         if e.When != "" {
//             upcoming = append(upcoming, e)
//         }
//     }
//     return upcoming
// }
//
// func (h *eventHandler) GetEventIndex(id string) (int, bool) {
//     for i, e := range h.Store {
//         if id == e.Id {
//             return i, true
//         }
//     }
//     return 0, false
// }
//
// func validEvent(reqEvent Event) (error, bool) {
//     // Making sure all fields are entered
//     missingFields := make([]string, 0)
//     if reqEvent.When == "" {
//         missingFields = append(missingFields, "when")
//     }
//     if reqEvent.Where == "" {
//         missingFields = append(missingFields, "where")
//     }
//     if reqEvent.What == "" {
//         missingFields = append(missingFields, "what")
//     }
//     if len(missingFields) != 0 {
//         errString := "missing json fields: "
//         for _, v := range missingFields {
//             errString += v + " "
//         }
//         return newHTTPError(nil, errString, http.StatusBadRequest), false
//     }
//
//     // Test whether "when" is in the correct format
//     _, err := time.Parse("15:04 02-01-06", reqEvent.When)
//     if err != nil {
//         return newHTTPError(
//             err,
//             "time not in '15:04 02-01-06' format",
//             http.StatusBadRequest,
//         ),
//         false
//     }
//     return nil, true
// }
