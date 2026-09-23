package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type dataResponse struct {
	Data any `json:"data"`
}

type errorResponse = Error

// RespondJSON writes data as a JSON response with the given status code,
// wrapped in the { "data": any } convention.
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dataResponse{Data: data})
}

// RespondError writes a JSON error response with the given status code,
// wrapped in the { "code": string, "message": string } convention.
func RespondError(w http.ResponseWriter, status int, code, message string) {
	RespondErrorWithFields(w, status, code, message, nil)
}

// RespondErrorWithFields writes a JSON error response with optional field errors.
func RespondErrorWithFields(
	w http.ResponseWriter,
	status int,
	code, message string,
	fields []FieldError,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Code:    code,
		Message: message,
		Fields:  fields,
	})
}

// RespondInvalidBody writes an invalid request body error response.
func RespondInvalidBody(w http.ResponseWriter, err error) {
	var requestError *Error
	if errors.As(err, &requestError) {
		status := requestError.Status
		if status == 0 {
			status = http.StatusBadRequest
		}

		code := requestError.Code
		if code == "" {
			code = "invalid_body"
		}

		RespondErrorWithFields(w, status, code, requestError.Message, requestError.Fields)
		return
	}

	RespondError(w, http.StatusBadRequest, "invalid_body", err.Error())
}

// RespondValidationError writes a validation error response.
func RespondValidationError(w http.ResponseWriter, err error) {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) || len(validationErrors) == 0 {
		RespondInvalidBody(w, err)
		return
	}

	fields := make([]FieldError, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		field := strings.ToLower(validationError.Field())
		fields = append(fields, FieldError{
			Field:   field,
			Rule:    validationError.Tag(),
			Param:   validationError.Param(),
			Message: validationMessage(validationError),
		})
	}

	firstField := fields[0]
	RespondErrorWithFields(
		w,
		http.StatusBadRequest,
		"invalid_"+firstField.Field,
		firstField.Field+" "+firstField.Message,
		fields,
	)
}
