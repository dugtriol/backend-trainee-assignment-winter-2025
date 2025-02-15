package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/service"
	"backend-trainee-assignment-winter-2025/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type transactionRoutes struct {
	transactionService service.Transaction
}

func newTransactionRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, transactionService service.Transaction,
) {
	u := transactionRoutes{transactionService: transactionService}
	route.Get("/sendCoin", u.sendCoin(ctx, log))
}

type inputSendCoin struct {
	ToUser string `json:"toUser" validate:"required,uuid"`
	Amount int    `json:"amount" validate:"required,numeric,gte=0"`
}

func (u *transactionRoutes) sendCoin(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		current := &entity.User{}
		if current, err = GetCurrentUserFromContext(r.Context()); err != nil {
			log.Info("transactionRoutes - GetCurrentUserFromContext", err)
			response.NewError(
				w,
				r,
				log,
				ErrNoUserInContext,
				http.StatusUnauthorized,
				MsgUserNotAuthorized,
			)
			return
		}
		var input inputSendCoin

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		transactionInput := service.TransactionInput{
			FromUserId: current.Id,
			ToUserId:   input.ToUser,
			Amount:     input.Amount,
		}

		if err = u.transactionService.Transfer(ctx, log, transactionInput); err != nil {
			if errors.Is(err, service.ErrSimilarId) {
				response.NewError(
					w,
					r,
					log,
					ErrSimilarId,
					http.StatusBadRequest,
					MsgSimilarId,
				)
				return
			}
			if errors.Is(err, service.ErrLowBalance) {
				response.NewError(
					w,
					r,
					log,
					ErrLowBalance,
					http.StatusBadRequest,
					MsgLowBalance,
				)
				return
			}
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}
	}
}
