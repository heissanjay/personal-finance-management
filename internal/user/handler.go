package user

import (
	"encoding/json"
	"net/http"

	"github.com/heissanjay/personal-finance-management/internal/auth"
	"github.com/heissanjay/personal-finance-management/internal/storage"
)

type Handler struct {
	AuthService *auth.AuthService
}

func NewHandler(authService *auth.AuthService) *Handler {
	return &Handler{AuthService: authService}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var u storage.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.AuthService.RegisterUser(u.Username, u.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully. Please login to continue"))
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u storage.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.LoginUser(u.Username, u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
