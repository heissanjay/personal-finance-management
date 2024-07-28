package storage

import (
	"context"
)

type DataStore interface {
	SaveUser(ctx context.Context, user User) error
	GetUserByUsername(ctx context.Context, username string) (User, error)

	SaveExpense(ctx context.Context, expense Expense, userID int) error
	UpdateExpense(ctx context.Context, expenseID int, expense Expense, userID int) error
	DeleteExpense(ctx context.Context, expenseID int, userID int) error
	ListExpenses(ctx context.Context, userID int) ([]Expense, error)
	GetExpenseByID(ctx context.Context, expenseID int, userID int) (Expense, error)
}

type User struct {
	ID       int
	Username string
	Password string
}

type Expense struct {
	ID     int     `json:"id"`
	UserID int     `json:"user_id"`
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
}
