package v1

import (
	"context"
	"log/slog"
	"net/http"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/service"
	"backend-trainee-assignment-winter-2025/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type infoRoutes struct {
	infoService        service.Info
	userService        service.User
	inventoryService   service.Inventory
	transactionService service.Transaction
}

func newInfoRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, infoService service.Info, userService service.User,
	inventoryService service.Inventory,
	transactionService service.Transaction,
) {
	u := infoRoutes{
		infoService: infoService, userService: userService, inventoryService: inventoryService,
		transactionService: transactionService,
	}
	route.Get("/info", u.info(ctx, log))
}

type outputInfo struct {
	Coins       int                     `json:"coins"`
	Inventory   []service.InfoInventory `json:"inventory"`
	CoinHistory service.CoinHistory     `json:"coinHistory"`
}

func (u *infoRoutes) info(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &entity.User{}
		var err error
		var inventories []entity.Inventory
		var transactions []entity.Transaction
		var infoInventory []service.InfoInventory
		var coinHistory service.CoinHistory

		if user, err = GetCurrentUserFromContext(r.Context()); err != nil {
			log.Info("infoRoutes - service.GetById", err)
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

		if inventories, err = u.inventoryService.GetByUserId(ctx, log, user.Id); err != nil {
			log.Info("infoRoutes - inventoryService.GetByUserId", err)
		}

		if transactions, err = u.transactionService.GetByUserId(ctx, log, user.Id); err != nil {
			log.Info("infoRoutes - transactionService.GetByUserId", err)
		}

		if infoInventory, coinHistory, err = u.infoService.Get(
			ctx,
			log,
			user.Id,
			inventories,
			transactions,
		); err != nil {
			log.Info("infoRoutes - len(infoInventory) == 0")
			response.NewError(
				w,
				r,
				log,
				ErrInternalServerErr,
				http.StatusInternalServerError,
				MsgInternalServerErr,
			)
			return
		}

		if len(infoInventory) == 0 {
			infoInventory = []service.InfoInventory{}
		}
		if len(coinHistory.Received) == 0 {
			coinHistory.Received = make([]service.TransactionReceived, 0)
		}
		if len(coinHistory.Sent) == 0 {
			coinHistory.Sent = make([]service.TransactionSent, 0)
		}

		output := outputInfo{
			Coins:       user.Amount,
			Inventory:   infoInventory,
			CoinHistory: coinHistory,
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, output)
	}
}
