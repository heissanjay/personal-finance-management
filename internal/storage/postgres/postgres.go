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
