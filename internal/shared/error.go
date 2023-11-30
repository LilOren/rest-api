package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/lil-oren/rest/internal/dto"
)

type (
	CustomErrorCode int
	CustomError     struct {
		Code            CustomErrorCode `json:"code"`
		Message         string          `json:"msg"`
		ResponseMessage string          `json:"message"`
	}
)

const (
	BadRequest CustomErrorCode = iota
	NotFound
	Forbidden
	Unauthorized
	InternalServer
)

func (ce CustomError) Error() string {
	return ce.ResponseMessage
}

func (ce CustomError) CreateHTTPErrorMessage() dto.JSONResponse {
	return dto.JSONResponse{
		Message: ce.ResponseMessage,
	}
}

func NewCustomError(code CustomErrorCode, message string) *CustomError {
	return &CustomError{
		Code:            code,
		ResponseMessage: message,
		Message:         strings.ToLower(message),
	}
}

func ValidationErrResponse(err error) string {
	var ve validator.ValidationErrors
	res := make(map[string]string)

	if !errors.As(err, &ve) {
		b, _ := json.Marshal(ve)
		return string(b)
	}

	for _, r := range ve {
		fieldName := strings.ToLower(r.Field())
		res[fieldName] = msgForTag(r)
	}

	b, _ := json.Marshal(res)
	return string(b)
}

func msgForTag(r validator.FieldError) string {
	switch r.Tag() {
	case "required":
		return fmt.Sprintf("%s field is a required field", r.Field())
	case "email":
		return fmt.Sprintf("%s field is not a valid email", r.Field())
	case "min":
		return fmt.Sprintf("%s field minimum is %s", r.Field(), r.Param())
	case "max":
		return fmt.Sprintf("%s field maximum is %s", r.Field(), r.Param())
	case "gte":
		return fmt.Sprintf("%s field less than %s", r.Field(), r.Param())
	case "lte":
		return fmt.Sprintf("%s field greater than %s", r.Field(), r.Param())
	case "gt":
		return fmt.Sprintf("%s field less than %s", r.Field(), r.Param())
	case "lt":
		return fmt.Sprintf("%s field greater than %s", r.Field(), r.Param())
	case "oneof":
		return fmt.Sprintf("%s field should be one of %s", r.Field(), strings.Join(strings.Split(r.Param(), " "), ","))
	case "alphanum":
		return fmt.Sprintf("%s field should only contains alphabet and/or number", r.Field())
	case "len":
		return fmt.Sprintf("%s field should be of length %s", r.Field(), r.Param())
	case "url":
		return fmt.Sprintf("%s field should be a URL", r.Field())
	case "e164":
		return fmt.Sprintf("%s field should be in +62<phone_number> format", r.Field())
	default:
		return ""
	}
}
