package data

import (
	"encoding/json"
	"io"
)

func (e Events) ToJSON(w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(e)
}

func (e *Event) ToJSON(w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(e)
}

func (e *Event) FromJSON(r io.Reader) error {
    d := json.NewDecoder(r)
    return d.Decode(e)
}

type JSONErrs []error

func (je JSONErrs) MarshalJSON() ([]byte, error) {
    res := make([]interface{}, len(je))
    for i, e := range je {
        if _, ok := e.(json.Marshaler); ok {
            res[i] = e
        } else {
            res[i] = e.Error()
        }
    }
    return json.Marshal(res)
}

func (je JSONErrs) ToJSON(w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(je)
}
