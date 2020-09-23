package tasks

import "fmt"

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

func wrongCommandError() error {
	return newError(4, "Command error")
}

func unsupportCommandError() error {
	return newError(5, "Command don't support right now")
}

func taskRunError(errMsg string) error {
	return newError(6, "Run task error: "+errMsg)
}

func scritpTaskError(err error) error {
	return newError(7, "Script task error: "+err.Error())
}

func uploadFileError(errMsg string) error {
	return newError(8, "Upload file error: "+errMsg)
}

func wrongArgsError(errMsg string) error {
	return newError(9, "Argument error: "+errMsg)
}
