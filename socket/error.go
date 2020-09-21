package socket

import (
	"fmt"
)

type socketError struct {
	statusCode int
	errMsg     string
}

func (e *socketError) Error() string {
	return fmt.Sprintf("Socket Error: %d - %s", e.statusCode, e.errMsg)
}

func newError(code int, msg string) error {
	return &socketError{statusCode: code, errMsg: msg}
}

func unsupportServerError() error {
	return newError(1, "Unsupport socket server type")
}

func serverError(serverType, msg string, err error) error {
	var errMsg string
	if msg != "" {
		errMsg += msg
	}
	if err != nil {
		errMsg += err.Error()
	}
	code := 2
	if serverType == "tcp" {
		code = 3
	}
	return newError(code, errMsg)
}
