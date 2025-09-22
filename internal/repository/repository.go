package repository

import (
	"context"
	"errors"

	"github.com/aibekfatkhulla/shop/internal/domain"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *repository {
	return &repository{conn: conn}
}

func (r *repository) CreateUser(ctx context.Context, user *domain.User) error {
	sqlStatement := `
		INSERT INTO users (id, name, password, email, number, address, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)
		RETURNING id
`
	err := r.conn.QueryRow(
		ctx,
		sqlStatement,
		user.ID,
		user.Name,
		user.Password,
		user.Email,
		user.Number,
		user.Address,
		user.Balance,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return err
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	sqlStatement :=
		`SELECT id, name, password, email, number, address, balance, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	var user domain.User
	err := r.conn.QueryRow(
		ctx,
		sqlStatement,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.Number,
		&user.Address,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrorUserNotFound
	}

	return &user, err
}
