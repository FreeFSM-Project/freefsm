package handlers

import (
	"net/http"

	"github.com/MartialM1nd/freefsm/internal/services"
	"github.com/MartialM1nd/freefsm/internal/templates"
)

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	templates.DashboardPage().Render(r.Context(), w)
}

func handleLogout(w http.ResponseWriter, r *http.Request, sessions *services.SessionService) {
	cookie, err := r.Cookie("session")
	if err == nil {
		sessions.Delete(r.Context(), cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session", Value: "", Path: "/", MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
