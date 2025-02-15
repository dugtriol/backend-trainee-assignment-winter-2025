package repoerrs

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrLowBalance    = errors.New("insufficient funds")
)
