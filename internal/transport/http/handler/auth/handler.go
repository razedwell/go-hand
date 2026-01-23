package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/razedwell/go-hand/internal/platform/logger"
	"github.com/razedwell/go-hand/internal/service/auth"
	"github.com/razedwell/go-hand/internal/service/user"
	"github.com/razedwell/go-hand/internal/transport/http/helpers"
	"github.com/razedwell/go-hand/internal/transport/http/middleware"
)

type Handler struct {
	userService *user.Service
	authService *auth.Service
	authMW      func(http.Handler) http.Handler
}

func NewHandler(userService *user.Service, authService *auth.Service, authMW func(http.Handler) http.Handler) *Handler {
	return &Handler{userService, authService, authMW}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /refresh", h.Refresh)

	// --- ADD THIS SECTION ---
	// Serves the index.html file at the route /test
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/transport/http/index.html")
	})
	// ------------------------

	protected := h.authMW

	mux.Handle("GET /logout", protected(http.HandlerFunc(h.Logout)))
	mux.Handle("GET /", protected(http.HandlerFunc(h.Home)))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   24 * 60 * 60,
	})

	helpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   accessToken,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	accessTokenVal := r.Context().Value(middleware.AcsKey)
	accessToken, ok := accessTokenVal.(string)
	if !ok {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		logger.Log.Printf("Invalid access token in context: %v", accessTokenVal)
		return
	}
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		cookie = &http.Cookie{}
	}

	refreshToken := cookie.Value

	if err := h.authService.Logout(r.Context(), accessToken, refreshToken); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh", // Must match the path used in Login
		HttpOnly: true,
		MaxAge:   -1, // Tells browser to delete immediately
		Expires:  time.Unix(0, 0),
	})

	helpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Logout successful",
	})
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.userService.RegisterUser(r.Context(), req); err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
	})
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Born&Razed",
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not provided", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	newAccessToken, err := h.authService.RefreshAccessToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Failed to refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token": newAccessToken,
	})
}
