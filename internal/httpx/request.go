package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

const maxBodySize = 1 << 20 // 1 MB

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Make validator errors use the JSON field name rather than
	// the Go struct field name.
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "" || name == "-" {
			return field.Name
		}

		return name
	})

	return v
}

// FieldError describes a validation error for a request field.
type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message"`
}

// Error describes a request error and its HTTP response details.
type Error struct {
	Status  int          `json:"-"`
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}

// Error returns the public error message.
func (e *Error) Error() string {
	return e.Message
}

// DecodeAndValidate decodes one JSON value from the HTTP body and validates it.
func DecodeAndValidate[T any](
	w http.ResponseWriter,
	r *http.Request,
	dst *T,
) *Error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode request body.
	if err := dec.Decode(dst); err != nil {
		return decodeError(err)
	}

	// Reject multiple JSON values:
	//
	// {"name":"John"} {"name":"Jane"}
	var extra json.RawMessage

	if err := dec.Decode(&extra); err != io.EOF {
		return &Error{
			Status:  http.StatusBadRequest,
			Code:    "invalid_body",
			Message: "request body must contain exactly one JSON value",
		}
	}

	// Validate decoded struct.
	if err := validate.StructCtx(r.Context(), dst); err != nil {
		var validationErrors validator.ValidationErrors

		if !errors.As(err, &validationErrors) {
			return &Error{
				Status:  http.StatusInternalServerError,
				Code:    "internal_error",
				Message: "failed to validate request",
			}
		}

		fields := make([]FieldError, 0, len(validationErrors))

		for _, e := range validationErrors {
			fields = append(fields, FieldError{
				Field:   strings.ToLower(e.Field()),
				Rule:    e.Tag(),
				Param:   e.Param(),
				Message: validationMessage(e),
			})
		}

		firstField := fields[0]
		return &Error{
			Status:  http.StatusBadRequest,
			Code:    "invalid_" + firstField.Field,
			Message: firstField.Field + " " + firstField.Message,
			Fields:  fields,
		}
	}

	return nil
}

func decodeError(err error) *Error {
	var syntaxError *json.SyntaxError
	var typeError *json.UnmarshalTypeError
	var maxBytesError *http.MaxBytesError

	switch {
	case errors.Is(err, io.EOF):
		return &Error{
			Status:  http.StatusBadRequest,
			Code:    "invalid_body",
			Message: "request body must not be empty",
		}

	case errors.As(err, &syntaxError):
		return &Error{
			Status: http.StatusBadRequest,
			Code:   "invalid_body",
			Message: fmt.Sprintf(
				"malformed JSON near position %d",
				syntaxError.Offset,
			),
		}

	case errors.As(err, &typeError):
		if typeError.Field != "" {
			return &Error{
				Status: http.StatusBadRequest,
				Code:   "invalid_body",
				Message: fmt.Sprintf(
					"invalid value for field %q",
					typeError.Field,
				),
			}
		}

		return &Error{
			Status:  http.StatusBadRequest,
			Code:    "invalid_body",
			Message: "JSON contains an invalid value type",
		}

	case errors.As(err, &maxBytesError):
		return &Error{
			Status: http.StatusRequestEntityTooLarge,
			Code:   "invalid_body",
			Message: fmt.Sprintf(
				"request body must not exceed %d bytes",
				maxBytesError.Limit,
			),
		}

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		field := strings.TrimPrefix(err.Error(), "json: unknown field ")

		return &Error{
			Status: http.StatusBadRequest,
			Code:   "invalid_body",
			Message: fmt.Sprintf(
				"unknown field %s",
				field,
			),
		}

	default:
		return &Error{
			Status:  http.StatusBadRequest,
			Code:    "invalid_body",
			Message: "invalid JSON request body",
		}
	}
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"

	case "email":
		return "must be a valid email address"

	case "min":
		return fmt.Sprintf("must satisfy minimum %s", e.Param())

	case "max":
		return fmt.Sprintf("must satisfy maximum %s", e.Param())

	case "len":
		return fmt.Sprintf("must have length %s", e.Param())

	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())

	default:
		if e.Param() != "" {
			return fmt.Sprintf(
				"failed validation %q with parameter %q",
				e.Tag(),
				e.Param(),
			)
		}

		return fmt.Sprintf("failed validation %q", e.Tag())
	}
}
