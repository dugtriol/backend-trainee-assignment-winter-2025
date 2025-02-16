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
	"github.com/go-playground/validator/v10"
)

type inventoryRoutes struct {
	inventoryService service.Inventory
}

func newInventoryRoutes(ctx context.Context, log *slog.Logger, route chi.Router, inventoryService service.Inventory) {
	u := inventoryRoutes{inventoryService: inventoryService}
	route.Get("/buy/{item}", u.buy(ctx, log))
}

type inputInventoryBuy struct {
	Item string `validate:"required"`
}

func (u *inventoryRoutes) buy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var current *entity.User
		if current, err = GetCurrentUserFromContext(r.Context()); err != nil {
			log.Info("inventoryRoutes - service.GetByID", "error", err.Error())
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
		if err = u.inventoryService.GetItem(ctx, log, current.ID, item); err != nil {
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
