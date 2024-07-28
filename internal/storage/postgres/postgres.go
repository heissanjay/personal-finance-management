package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/heissanjay/personal-finance-management/internal/storage"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{DB: db}
}

func (p *PostgresDB) InitDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(60) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY, 
		user_id INT NOT NULL,
		title VARCHAR(255) NOT NULL,
		amount NUMERIC(10,2) NOT NULL,
		date DATE NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err := p.DB.Exec(query)
	return err
}

func (p *PostgresDB) SaveUser(ctx context.Context, user storage.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := p.DB.ExecContext(ctx, query, user.Username, user.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.New("user already exists")
			}
		}
		return err
	}
	return nil
}

func (p *PostgresDB) GetUserByUsername(ctx context.Context, username string) (storage.User, error) {
	var user storage.User
	query := `SELECT id, username, password FROM users WHERE username = $1`
	err := p.DB.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

func (p *PostgresDB) SaveExpense(ctx context.Context, expense storage.Expense, userID int) error {
	query := `INSERT INTO expenses (user_id, title, amount, date) VALUES  ($1, $2, $3, $4)`
	_, err := p.DB.ExecContext(ctx, query, userID, expense.Title, expense.Amount, expense.Date)
	return err
}

func (p *PostgresDB) UpdateExpense(ctx context.Context, expenseID int, expense storage.Expense, userID int) error {
	query := `UPDATE expenses SET title = $1, amount = $2, date = $3 where id = $4 AND user_id = $5`
	_, err := p.DB.ExecContext(ctx, query, expense.Title, expense.Amount, expense.Date, expenseID, userID)
	return err
}

func (p *PostgresDB) DeleteExpense(ctx context.Context, expenseID int, userID int) error {
	query := `DELETE FROM expenses WHERE id = $1 AND user_id = $2`
	_, err := p.DB.ExecContext(ctx, query, expenseID, userID)
	return err
}

func (p *PostgresDB) ListExpenses(ctx context.Context, userID int) ([]storage.Expense, error) {
	query := `SELECT * FROM expenses WHERE user_id = $1`
	rows, err := p.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []storage.Expense

	for rows.Next() {
		var expense storage.Expense
		if err := rows.Scan(&expense.ID, &expense.UserID, &expense.Title, &expense.Amount, &expense.Date); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (p *PostgresDB) GetExpenseByID(ctx context.Context, expenseID int, userID int) (storage.Expense, error) {
	var expense storage.Expense
	query := `SELECT id, user_id, title, amount, date FROM expenses WHERE id = $1 AND user_id = $2`
	err := p.DB.QueryRowContext(ctx, query, expenseID, userID).Scan(&expense.ID, &expense.UserID, &expense.Title, &expense.Amount, &expense.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return expense, errors.New("expense not found")
		}
		return expense, err
	}
	return expense, nil
}
