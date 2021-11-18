package main

import(
    "net/http"
    "encoding/json"
    "log"
)

// Wrapper around handlers that deals with errors
type rootHandler func(http.ResponseWriter, *http.Request) error

// In order to be used with http.Handle it needs to implement serveHTTP method
func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    err := fn(w, r)
    if err == nil {
        return
    }
    // Error handling logic starts here:
    log.Printf("An error occured: %v", err)

    // Check if it is a ClientError
    clientError, ok := err.(ClientError)
    if !ok {
        w.WriteHeader(500)
        return
    }
    body, err := clientError.ResponseBody()
    if err != nil {
        log.Printf("An error occured: %v", err)
        w.WriteHeader(500)
        return
    }
    status, headers := clientError.ResponseHeaders()
    for k, v := range headers {
        w.Header().Set(k, v)
    }
    w.WriteHeader(status)
    w.Write(body)

}

func (h *eventHandler) EventsHandler(w http.ResponseWriter, r *http.Request) error {
    switch r.Method {
    case http.MethodGet:
        return h.Index(w, r)
    case http.MethodPost:
        return h.Create(w, r)
    default:
        return newHTTPError(
            nil,
            "method not allowed",
            http.StatusMethodNotAllowed,
        )
    }
}

func (h *eventHandler) Index(w http.ResponseWriter, r *http.Request) error {
    h.Lock()
    defer h.Unlock()

    jsonData, err := json.Marshal(h.Store)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return newHTTPError(
            err,
            "error fetching event data.",
            http.StatusInternalServerError,
        )
    }

    w.Header().Add("content-type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
    return nil
}

func (h *eventHandler) Create(w http.ResponseWriter, r *http.Request) error {
    w.Write([]byte("hello lovely :)"))
    return nil
}

func main() {
    h := newEventHandler()

    // Converting our eventHandler methods to rootHandler functions
    // to separate error handling with HTTP handling
    http.Handle("/events", rootHandler(h.EventsHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
