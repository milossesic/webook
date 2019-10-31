package util

import (
	"fmt"
	"net/http"
)

// AppError is error struct
type AppError struct {
	Err      string `json:"msg"`
	HTTPCode int    `json:"code"`
	Cause    string `json:"cause"`
}

// Error returns a JSON base reponse that contains error msg and status code
func Error(w http.ResponseWriter, reqID string, httpCode int, cause string) {

	var err = AppError{
		Err:      fmt.Sprintf("%s, Request ID: %s", http.StatusText(httpCode), reqID),
		HTTPCode: httpCode,
		Cause:    cause,
	}

	JSON(w, err, httpCode)
}
