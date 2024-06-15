package errorhandler

import "net/http"

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Err  error
	Code int
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

func GetStatusCode(err error) int {
	if err != nil {
		switch e := err.(type) {
		case Error:
			return e.Status()
		default:
			return http.StatusInternalServerError
		}
	}
	return http.StatusOK
}
