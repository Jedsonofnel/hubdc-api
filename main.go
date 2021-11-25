package main

import(
    "net/http"
    "encoding/json"
    "log"
    "io/ioutil"
    "strconv"
    "strings"
    "fmt"
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
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    body, err := clientError.ResponseBody()
    if err != nil {
        log.Printf("An error occured: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    status, headers := clientError.ResponseHeaders()
    for k, v := range headers {
        w.Header().Set(k, v)
    }
    w.WriteHeader(status)
    w.Write(body)
}

func (h *eventHandler) Events(w http.ResponseWriter, r *http.Request) error {
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
        return newHTTPError(
            err,
            "error fetching event data",
            http.StatusInternalServerError,
        )
    }

    w.Header().Add("content-type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
    return nil
}

func (h *eventHandler) Create(w http.ResponseWriter, r *http.Request) error {
    user, pass, ok := r.BasicAuth()
    if !ok || user != "hubdc-admin" || pass != h.Password {
        return newHTTPError(nil, "invalid authorisation", http.StatusUnauthorized)
    }

    bodyBytes, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        return newHTTPError(err, "error reading request", http.StatusInternalServerError)
    }

    ct := r.Header.Get("content-type")
    if ct != "application/json" {
        return newHTTPError(err, "need content-type: application/json", http.StatusBadRequest)
    }

    var reqEvent Event
    err = json.Unmarshal(bodyBytes, &reqEvent)
    if err != nil {
        return newHTTPError(err, "error parsing json", http.StatusBadRequest)
    }

    if err, ok := validEvent(reqEvent); !ok {
        return err
    }

    reqEvent.Id = strconv.Itoa(len(h.Store))

    // Add good data to the store
    h.Lock()
    defer h.Unlock()
    h.Store = append(h.Store, reqEvent)

    // Serialise Baby
    err = h.SerialiseBaby()
    if err != nil {
        return err
    }
    w.WriteHeader(http.StatusOK)
    return nil
}

func (h *eventHandler) Event(w http.ResponseWriter, r *http.Request) error {
    switch r.Method {
    case http.MethodGet:
        return h.Show(w, r, event)
    case http.MethodPut:
        return h.Update(w, r, event)
    default:
        return newHTTPError(
            nil,
            "method not allowed",
            http.StatusMethodNotAllowed,
        )
    }
    // Get requested id from url parameter
    parts := strings.Split(r.URL.String(), "/")
    if len(parts) != 3 {
        return newHTTPError(nil, "invalid url", http.StatusNotFound)
    }
    id := parts[2]

    h.Lock()
    event, ok := h.FindWithID(id)
    h.Unlock()
    if !ok {
        return newHTTPError(
            nil,
            fmt.Sprintf("event '%v' not found", id),
            http.StatusNotFound,
        )
    }

    switch r.Method {
    case http.MethodGet:
        return h.Show(w, r, event)
    case http.MethodPut:
        return h.Update(w, r, event)
    default:
        return newHTTPError(
            nil,
            "method not allowed",
            http.StatusMethodNotAllowed,
        )
    }
}

func (h *eventHandler) Show(w http.ResponseWriter, r *http.Request, e Event) error {
    jsonData, err := json.Marshal(e)
    if err != nil {
        return newHTTPError(
            err,
            "error parsing event data",
            http.StatusInternalServerError,
        )
    }

    w.Header().Add("content-type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
    return nil
}

func (h *eventHandler) Update(w http.ResponseWriter, r *http.Request, e Event) error {
    user, pass, ok := r.BasicAuth()
    if !ok || user != "hubdc-admin" || pass != h.Password {
        return newHTTPError(nil, "invalid authorisation", http.StatusUnauthorized)
    }

    bodyBytes, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        return newHTTPError(err, "error reading request", http.StatusInternalServerError)
    }

    ct := r.Header.Get("content-type")
    if ct != "application/json" {
        return newHTTPError(err, "need content-type: application/json", http.StatusBadRequest)
    }

    var reqEvent Event
    err = json.Unmarshal(bodyBytes, &reqEvent)
    if err != nil {
        return newHTTPError(err, "error parsing json", http.StatusBadRequest)
    }

    if err, ok := validEvent(reqEvent); !ok {
        return err
    }

    intId, _ := strconv.Atoi(e.Id)

    h.Lock()
    h.Store[intId].When = reqEvent.When
    h.Store[intId].Where = reqEvent.Where
    h.Store[intId].What = reqEvent.What
    defer h.Unlock()

    err = h.SerialiseBaby()
    if err != nil {
        return err
    }
    w.WriteHeader(http.StatusOK)
    return nil
}

func (h *eventHandler) ServeUpcoming(w http.ResponseWriter, r *http.Request) error {
    if r.Method != http.MethodGet {
        return newHTTPError(nil, "method not allowed", http.StatusMethodNotAllowed)
    }

    jsonData, err := json.Marshal(h.Upcoming())
    if err != nil {
        return newHTTPError(err, "error fetching event data", http.StatusInternalServerError)
    }

    w.Header().Add("content-type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
    return nil
}

func main() {
    h := newEventHandler()
    // Index and create routes
    http.Handle("/events", rootHandler(h.Events))
    // Show, update and delete routes
    http.Handle("/event/", rootHandler(h.Event))
    // Helper route for getting array of next three events
    http.Handle("/events/upcoming", rootHandler(h.ServeUpcoming))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
