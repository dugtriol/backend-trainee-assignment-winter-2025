package service

import "fmt"

var (
	ErrCannotHashPassword = fmt.Errorf("cannot hash password")

	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")

	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrInvalidPassword   = fmt.Errorf("invalid password")
	ErrCannotUpdateUser  = fmt.Errorf("cannot update user")

	ErrInvalidMerchType = fmt.Errorf("invalid merch type")

	ErrAlreadyExists = fmt.Errorf("already exists")
	ErrCannotCreate  = fmt.Errorf("cannot create")
	ErrCannotGet     = fmt.Errorf("cannot get")
)
