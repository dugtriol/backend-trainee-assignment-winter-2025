package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/service"
	mw "backend-trainee-assignment-winter-2025/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const api = "/api"

func NewRouter(ctx context.Context, log *slog.Logger, route *chi.Mux, services *service.Services) {
	route.Use(middleware.Logger)
	route.Use(middleware.RequestID)
	route.Use(middleware.Recoverer)
	route.Use(middleware.URLFormat)
	route.Use(mw.New(log))
	route.Use(render.SetContentType(render.ContentTypeJSON))

	route.Route(
		api, func(r chi.Router) {
			newUserRoutes(ctx, log, r, services.User)
			r.Group(
				func(g chi.Router) {
					g.Use(AuthMiddleware(ctx, log, services.User))
					g.Get("/ping", Ping())
				},
			)
		},
	)
}

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output := &entity.User{}
		var err error
		if output, err = GetCurrentUserFromCTX(r.Context()); err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(fmt.Sprintf("id - %s, name - %s", output.Id, output.Username)))
		if err != nil {
			return
		}
	}
}
