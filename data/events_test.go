package data

import "testing"

func TestChecksValidation(t *testing.T) {
    e := &Event{
        What: "Normal Hub Sesion",
        Where: "HPH",
        When: "23:15 12-03-22",
    }

    err := e.Validate()

    if err != nil {
        t.Fatal(err)
    }
}
