package auth

import (
	"encoding/json"
	"net/http"

	"github.com/razedwell/go-hand/internal/service/auth"
	"github.com/razedwell/go-hand/internal/service/user"
	"github.com/razedwell/go-hand/internal/transport/http/helpers"
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

	protected := h.authMW

	mux.Handle("GET /", protected(http.HandlerFunc(h.Home)))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "Incorrect login credentials", http.StatusUnauthorized)
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"data":    user,
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
