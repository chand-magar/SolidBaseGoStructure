package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg string

	firstErr := errs[0]
	switch firstErr.ActualTag() {
	case "required":
		errMsg = fmt.Sprintf("field %s is required", firstErr.Field())
	default:
		errMsg = fmt.Sprintf("field %s is invalid", firstErr.Field())
	}

	return Response{
		Status: StatusError,
		Error:  errMsg,
	}
}
