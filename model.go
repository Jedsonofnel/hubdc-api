package main

import (
    "os"
    "encoding/json"
    "log"
    "io/ioutil"
    "sync"
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
