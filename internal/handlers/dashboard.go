package handlers

import (
	"net/http"
	"time"

	"github.com/MartialM1nd/freefsm/internal/middleware"
	"github.com/MartialM1nd/freefsm/internal/services"
	"github.com/MartialM1nd/freefsm/internal/templates"
)

type DashboardHandler struct {
	dashboardSvc *services.DashboardService
	timeEntrySvc *services.TimeEntryService
}

func NewDashboardHandler(dashboardSvc *services.DashboardService, timeEntrySvc *services.TimeEntryService) *DashboardHandler {
	return &DashboardHandler{dashboardSvc: dashboardSvc, timeEntrySvc: timeEntrySvc}
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	loc := middleware.CompanyLocation(r.Context())
	clockWidget := templates.ClockWidgetData{}
	user, _ := middleware.UserFromContext(r.Context())
	stats := services.DashboardStats{}
	if isAdminOrDispatcher(user) {
		stats, _ = h.dashboardSvc.Stats(r.Context(), loc)
	}
	if user != nil {
		activeEntry, err := h.timeEntrySvc.GetActiveByUser(r.Context(), user.ID)
		if err == nil && activeEntry != nil {
			duration := services.TimeEntryDuration(activeEntry.ClockIn, time.Time{})
			clockWidget = templates.ClockWidgetData{
				IsClockedIn: true,
				Duration:    duration,
				ClockInTime: activeEntry.ClockIn.In(loc).Format("Jan 2, 2006 3:04 PM"),
			}
		}
	}

	templates.DashboardPage(templates.DashboardData{
		Stats:       stats,
		ClockWidget: clockWidget,
	}).Render(r.Context(), w)
}

func handleLogout(w http.ResponseWriter, r *http.Request, sessions *services.SessionService, activitySvc *services.ActivityService) {
	u, _ := middleware.UserFromContext(r.Context())
	if u != nil && activitySvc != nil {
		activitySvc.Record(r.Context(), u.ID, "logged_out", "user", u.ID, map[string]interface{}{
			"entity_name": u.Name,
			"actor_name":  u.Name,
		})
	}
	cookie, err := r.Cookie("session")
	if err == nil {
		sessions.Delete(r.Context(), cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session", Value: "", Path: "/", MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
