package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ParseAndValidate parses the JSON request body into the given struct pointer v and validates it.
func ParseAndValidate(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json payload: %w", err)
	}

	if err := validate.Struct(v); err != nil {
		// Map validation errors to a simpler format
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, fieldError := range validationErrors {
				errMsgs = append(errMsgs, fmt.Sprintf("field '%s' failed on the '%s' tag", fieldError.Field(), fieldError.Tag()))
			}
			return fmt.Errorf("validation failed: %v", errMsgs)
		}
		return err
	}

	return nil
}
