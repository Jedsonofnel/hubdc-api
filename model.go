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
    storeFile   string  `json:"-"`
}

func newEventHandler() *eventHandler {
    var e eventHandler
    e.storeFile = "eventStore.json"

    // Open or create json "storeFile" if non-existent
    f, err := os.OpenFile(e.storeFile, os.O_CREATE, 0644)
    if err != nil {
        log.Println(err)
    }
    defer f.Close()

    byteValue, err := ioutil.ReadAll(f)
    json.Unmarshal(byteValue, &e)

    return &e
}
