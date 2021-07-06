package postgres

import (
	"errors"
)

var (
	NotFoundError      = errors.New("movie not found")
	AlreadyExists      = errors.New("already exists")
	UnmarshallError    = errors.New("movie json unmarshalling error")
	InvalidFilterError = errors.New("invalid filter")
	InvalidVoteError   = errors.New("invalid vote")
	RatingUpdateError  = errors.New("can't update rating")
	ViewUpdateError    = errors.New("can't update rating")
	BadParamsError     = errors.New("invalid vote")
	InvalidViewCheck   = errors.New("invalid view check")
	InvalidViewAdd     = errors.New("invalid view add")
)
