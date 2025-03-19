package rest_err

import (
	"net/http"
)

type RestErr struct {
	Message string `json:"message"`
	Err string `json:"error"`
	Code int	`json:"code"`
	Causes []Causes `json:"causes,omitempty"`
}

type Causes struct {
	Field string `json:"field"`
	Message string `json:"message"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err: "Bad_request",
		Code: http.StatusBadRequest,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err: "Internal server error",
		Code: http.StatusInternalServerError,
	}
}


func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err: "Not Found",
		Code: http.StatusNotFound,
	}
}