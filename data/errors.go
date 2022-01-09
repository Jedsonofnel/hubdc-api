package data

type JSONError struct {
    Message string `json:"message"`
}

func NewJE(msg string) JSONError {
    return JSONError{Message: msg}
}

type JSONErrors struct {
    Errors []JSONError `json:"errors"`
}

func NewJEs(msgs ...string) JSONErrors {
    var jes JSONErrors
    for _, msg := range msgs {
        jes.Errors = append(jes.Errors, NewJE(msg))
    }
    return jes
}
