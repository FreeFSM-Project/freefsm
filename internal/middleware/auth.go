package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MartialM1nd/freefsm/internal/services"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(sessions *services.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			userID, err := sessions.Validate(r.Context(), cookie.Value)
			if err != nil {
				http.SetCookie(w, &http.Cookie{
					Name: "session", Value: "", Path: "/", MaxAge: -1,
				})
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIDKey).(int64)
	return id, ok
}

func Public(path string) bool {
	return path == "/login" || path == "/setup" || path == "/health" || strings.HasPrefix(path, "/static/")
}
