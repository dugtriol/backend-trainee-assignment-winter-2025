package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"backend-trainee-assignment-winter-2025/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(ctx context.Context, log *slog.Logger, route chi.Router, userService service.User) {
	u := userRoutes{userService: userService}
	route.Route(
		"/", func(r chi.Router) {
			r.Post("/auth", u.auth(ctx, log))
		},
	)
}

type inputAuth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *userRoutes) auth(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputAuth
		var err error
		var tokenString string

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			newErrorResponse(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			newErrorValidateResponse(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}
		if tokenString, err = u.userService.Auth(
			ctx, log, service.AuthInput{
				Username: input.Username,
				Password: input.Password,
			},
		); err != nil {
			if errors.Is(err, service.ErrInvalidPassword) {
				newErrorResponse(w, r, log, err, http.StatusBadRequest, MsgInvalidPasswordErr)
				return
			}
			log.Error(err.Error())
			newErrorResponse(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		type response struct {
			Token string `json:"token"`
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response{Token: tokenString})
	}
}
