package util

import (
	"encoding/json"
	"net/http"
)

// JSON method returns json based response
func JSON(w http.ResponseWriter, val interface{}, code int) {

	var (
		b   []byte
		err error
	)

	if b, err = json.Marshal(val); err != nil {
		JSON(w, &AppError{Err: http.StatusText(500), HTTPCode: 500, Cause: "Json marshalling failed"}, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}
