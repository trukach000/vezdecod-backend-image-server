package httpext

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrBadRequest(w http.ResponseWriter) {
	Abort(w, errors.New("invalid Request"), http.StatusBadRequest)
}

func ErrInternalServerError(w http.ResponseWriter) {
	Abort(w, errors.New("internal Server Error"), http.StatusInternalServerError)
}

func ErrForbidden(w http.ResponseWriter) {
	Abort(w, errors.New("forbidden"), http.StatusForbidden)
}
