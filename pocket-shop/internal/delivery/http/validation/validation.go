package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"pocket-shop/internal/delivery/http/response"
)

func ErrorMessages(err error) []string {
	if err == nil {
		return nil
	}
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) || len(errs) == 0 {
		return []string{err.Error()}
	}
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, fmt.Sprintf("%s: %s", strings.ToLower(e.Field()), e.Tag()))
	}
	return msgs
}

func ErrorResponse(err error) response.ErrorResponse {
	msgs := ErrorMessages(err)
	if len(msgs) == 0 {
		return response.ErrorResponse{Error: "bad_request", Message: "validation failed"}
	}
	return response.ErrorResponse{
		Error:   "bad_request",
		Message: msgs[0],
		Errors:  msgs,
	}
}
