package taskerclient

import (
	"encoding/json"
)

type Error struct {
	Error string `json:"error"`
}

func NewError(err error) *Error {
	return &Error{Error: err.Error()}
}

func (e *Error) Marshall() ([]byte, error) {
	return json.Marshal(e)
}