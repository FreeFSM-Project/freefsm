package handlers

import (
	"net/http"

	"github.com/MartialM1nd/freefsm/internal/services"
	"github.com/MartialM1nd/freefsm/internal/templates"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db       *pgxpool.Pool
	sessions *services.SessionService
}

func NewAuthHandler(db *pgxpool.Pool, sessions *services.SessionService) *AuthHandler {
	return &AuthHandler{db: db, sessions: sessions}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.showLogin(w, r)
	case http.MethodPost:
		h.login(w, r)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (h *AuthHandler) showLogin(w http.ResponseWriter, r *http.Request) {
	if needsSetup(r.Context(), h.db) {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}
	templates.LoginPage(templates.LoginPageData{
		Error: r.URL.Query().Get("error"),
	}).Render(r.Context(), w)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error=invalid+form", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	var id int64
	var hash string
	err := h.db.QueryRow(r.Context(),
		`SELECT id, password_hash FROM users WHERE email = $1 AND is_active = true`, email,
	).Scan(&id, &hash)
	if err != nil {
		http.Redirect(w, r, "/login?error=invalid+credentials", http.StatusSeeOther)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		http.Redirect(w, r, "/login?error=invalid+credentials", http.StatusSeeOther)
		return
	}

	token, err := h.sessions.Create(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session", Value: token, Path: "/",
		HttpOnly: true, SameSite: http.SameSiteLaxMode,
		MaxAge: 604800,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
