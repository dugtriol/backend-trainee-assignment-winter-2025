package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"backend-trainee-assignment-winter-2025/pkg/response"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	MsgInvalidReq         = "Invalid request"
	MsgFailedParsing      = "Failed to parse data"
	MsgInternalServerErr  = "Internal server error"
	MsgInvalidPasswordErr = "Invalid password"

	MsgUserNotFound = "User not found"

	MsgForbidden         = "Forbidden"
	MsgUserAlreadyExists = "User already exists"
)

func newErrorResponse(
	w http.ResponseWriter, r *http.Request, log *slog.Logger, err error, errStatus int, message string,
) {
	log.Error(message, err)
	w.WriteHeader(errStatus)
	render.JSON(w, r, response.MakeResponse(message))
}

func newErrorValidateResponse(
	w http.ResponseWriter, r *http.Request, log *slog.Logger, errStatus int, message string,
	err error,
) {
	var validateErr validator.ValidationErrors
	errors.As(err, &validateErr)

	log.Error(message, err)
	w.WriteHeader(errStatus)
	render.JSON(w, r, response.ValidationError(validateErr))
}
