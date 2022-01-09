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
