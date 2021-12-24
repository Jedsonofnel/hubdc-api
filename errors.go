package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type ClientError interface {
    Error() string
    ResponseBody() ([]byte, error)
    ResponseHeaders() (int, map[string]string)
}

type HTTPError struct {
    Cause   error   `json:"-"`
    Detail  string  `json:"detail"`
    Status  int     `json:"-"`
}

func (e *HTTPError) Error() string {
    if e.Cause == nil {
        return e.Detail
    }
    return e.Detail + " : " + e.Cause.Error()
}

func (e *HTTPError) ResponseBody() ([]byte, error) {
    body, err := json.Marshal(e)
    if err != nil {
        return nil, fmt.Errorf("Error while parsing response body: %v", err)
    }
    return body, nil
}

func (e *HTTPError) ResponseHeaders() (int, map[string]string) {
    return e.Status, map[string]string {
        "Content-type": "application/json; charset=utf-8",
    }
}

// Used to return
func newHTTPError(err error, detail string, status int) error {
    return &HTTPError{
        Cause:  err,
        Detail: detail,
        Status: status,
    }
}

func sqlError(e error) error {
    return &HTTPError{
        Cause: e,
        Detail: "error fetching event data",
        Status: http.StatusInternalServerError,
    }
}
