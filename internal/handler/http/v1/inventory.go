package v1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/service"
	"backend-trainee-assignment-winter-2025/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type inventoryRoutes struct {
	inventoryService service.Inventory
}

func newInventoryRoutes(ctx context.Context, log *slog.Logger, route chi.Router, inventoryService service.Inventory) {
	u := inventoryRoutes{inventoryService: inventoryService}
	route.Get("/buy/{item}", u.buy(ctx, log))
	route.Get("/pink-buy", u.PingBuy(ctx, log))
}

func (u *inventoryRoutes) PingBuy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output := &entity.User{}
		var err error
		if output, err = GetCurrentUserFromContext(r.Context()); err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		//w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(fmt.Sprintf("id - %s, name - %s", output.Id, output.Username)))
		if err != nil {
			return
		}
	}
}

type inputInventoryBuy struct {
	Item string `validate:"required"`
}

func (u *inventoryRoutes) buy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		current := &entity.User{}
		if current, err = GetCurrentUserFromContext(r.Context()); err != nil {
			log.Info("inventoryRoutes - service.GetById", "error", err.Error())
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
		item := chi.URLParam(r, "item")

		if err = validator.New().Struct(inputInventoryBuy{Item: item}); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}
		if err = u.inventoryService.GetItem(ctx, log, current.Id, item); err != nil {
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
