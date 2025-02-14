package v1

import "fmt"

var (
	MsgInvalidReq         = "Invalid request"
	MsgFailedParsing      = "Failed to parse data"
	MsgInternalServerErr  = "Internal server error"
	MsgInvalidPasswordErr = "Invalid password"

	MsgUserNotFound = "User not found"

	MsgForbidden         = "Forbidden"
	MsgUserAlreadyExists = "User already exists"

	ErrInvalidToken     = fmt.Errorf("invalid token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")
)
