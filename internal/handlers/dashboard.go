package handlers

import (
	"net/http"

	"github.com/MartialM1nd/freefsm/internal/services"
	"github.com/MartialM1nd/freefsm/internal/templates"
)

type DashboardHandler struct {
	dashboardSvc *services.DashboardService
}

func NewDashboardHandler(dashboardSvc *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardSvc: dashboardSvc}
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	stats, _ := h.dashboardSvc.Stats(r.Context())
	templates.DashboardPage(templates.DashboardData{
		Stats: stats,
	}).Render(r.Context(), w)
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
