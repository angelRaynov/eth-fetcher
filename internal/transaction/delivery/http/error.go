package http

import "errors"

type ErrorResponse struct {
	Message string `json:"message"`
}

var ErrNoRecords = errors.New("no records")
