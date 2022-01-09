package data

type JSONError struct {
    Message string `json:"message"`
}

func NewJE(msg string) JSONError {
    return JSONError{Message: msg}
}
