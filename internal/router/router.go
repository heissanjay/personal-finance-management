package router

import (
	"github.com/gorilla/mux"
	"github.com/heissanjay/personal-finance-management/internal/auth"
	"github.com/heissanjay/personal-finance-management/internal/expense"
	"github.com/heissanjay/personal-finance-management/internal/user"
)

func NewRouter(user *user.Handler, expenseHandler *expense.ExpenseHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", user.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", user.LoginHandler).Methods("POST")

	r.HandleFunc("/expenses", expenseHandler.CreateExpenseHandler).Methods("POST")
	r.HandleFunc("/expenses/{id}", expenseHandler.UpdateExpenseHandler).Methods("PUT")
	r.HandleFunc("/expenses/{id}", expenseHandler.DeleteExpenseHandler).Methods("DELETE")
	r.HandleFunc("/expenses/{id}", expenseHandler.GetExpenseByIdHandler).Methods("GET")
	r.HandleFunc("/expenses", expenseHandler.ListExpensesHandler).Methods("GET")

	r.Use(auth.Middleware)

	return r
}
