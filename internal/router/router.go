package router

import (
	"github.com/gorilla/mux"
	"github.com/heissanjay/personal-finance-management/internal/auth"
	"github.com/heissanjay/personal-finance-management/internal/user"
)

func NewRouter(user *user.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", user.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", user.LoginHandler).Methods("POST")
	r.Use(auth.Middleware)

	return r
}
