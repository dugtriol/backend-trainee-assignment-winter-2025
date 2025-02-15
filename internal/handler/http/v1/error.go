package v1

import "fmt"

var (
	MsgInvalidReq         = "Invalid request"
	MsgFailedParsing      = "Failed to parse data"
	MsgInternalServerErr  = "Internal server error"
	MsgInvalidPasswordErr = "Invalid password"

	MsgUserNotFound      = "User not found"
	MsgUserNotAuthorized = "User not authorized"

	MsgForbidden         = "Forbidden"
	MsgUserAlreadyExists = "User already exists"

	ErrInvalidToken     = fmt.Errorf("invalid token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")
	ErrUserGet          = fmt.Errorf("user not get from database")
	ErrNoUserInContext  = fmt.Errorf("no user in the context")
)
