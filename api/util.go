package api

import "net/http"

type ErrorResponse struct {
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{ErrorCode: http.StatusBadRequest, ErrorMessage: message}
}

func NewOKResponse() *ErrorResponse {
	return &ErrorResponse{ErrorCode: http.StatusOK, ErrorMessage: "OK"}
}
