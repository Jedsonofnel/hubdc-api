package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

func msgForTag(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return fmt.Sprintf("Field '%v' required", fe.Field())
    case "when":
        return "Time must be in format '15:04 02-01-06'"
    }
    return fe.Error()
}

func (e *Event) Validate() JSONErrs {
    validate := validator.New()
    validate.RegisterValidation("when", validateWhen)

    err := validate.Struct(e)

    if err == nil {
        return JSONErrs{}
    }


    var errs JSONErrs
    for _, err := range err.(validator.ValidationErrors) {
        errs = append(errs, errors.New(msgForTag(err)))
    }
    return errs
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
