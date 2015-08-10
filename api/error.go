package api

import "log"

type Error struct {
	Error string `json:"error"`
}

func NewError(msg string) *Error {
	return &Error{Error: msg}
}

func CheckErr(err error, msg string) *Error {
	if err != nil {
		log.Fatalln(msg, err)
		return NewError(msg)
	}
	return nil
}
