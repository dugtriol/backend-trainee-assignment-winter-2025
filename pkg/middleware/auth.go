package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"backend-trainee-assignment-winter-2025/internal/entity"
)

const CurrentUserKey = "currentUser"

func AuthMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			//r = r.WithContext(context.WithValue(r.Context(), key, val))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func GetCurrentUserFromCTX(ctx context.Context) (*entity.User, error) {
	errNoUserInContext := errors.New("no user in context")
	log.Println(fmt.Sprintf("ctx.Value(CurrentUserKey): %s", ctx.Value(CurrentUserKey)))
	if ctx.Value(CurrentUserKey) == nil {
		return nil, errNoUserInContext
	}

	user, ok := ctx.Value(CurrentUserKey).(entity.User)
	//log.Info(fmt.Sprintf("ctx.Value(CurrentUserKey).(*entity.User) %v",user))
	if !ok || user.Id == "" {
		return nil, errNoUserInContext
	}

	return &user, nil
}
