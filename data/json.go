package data

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
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

func (e *Event) Validate() error {
    validate := validator.New()
    validate.RegisterValidation("when", validateWhen)
    return validate.Struct(e)
}

func validateWhen(fl validator.FieldLevel) bool {
    // when is of format 15:04 02-01-06
    whenFmt := "15:04 02-01-06"
    _, err := time.Parse(whenFmt, fl.Field().String())

    if err != nil {
        return false
    }

    return true
}
