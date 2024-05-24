package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
}

func Error(w http.ResponseWriter, statusCode int, args ...interface{}) {
	message := http.StatusText(statusCode) // default message

	// determine if an error or string arg was passed in
	// set the message accordingly
	if len(args) != 0 {
		switch v := args[0].(type) {
		case string:
			message = v
		case error:
			message = v.Error()
		}
	}

	errString, _ := json.Marshal(ErrorMessage{
		Message:    message,
		StatusCode: statusCode,
	})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	w.Write(errString)
}
