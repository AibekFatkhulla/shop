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

func (r *repository) NewService(ctx context.Context, user *domain.User) error {
	//TODO implement me

	panic("implement me")
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

func (r *repository) UpdateUser(ctx context.Context, user *domain.User) error {
	sqlStatement :=
		`UPDATE users
		SET name = $2,
		    password = $3,
		    email = $4,
		    number =$5,
		    address = $6,
		    updated_at = $7,
		WHERE id = $1
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
		user.UpdatedAt,
	).Scan(&user.ID)
	return err
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, name, password, email, number, address, balance, created_at, updated_at
		FROM users
		WHERE id = $1
		`

	user := &domain.User{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
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
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrorUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *repository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, name, email, password, number, address, balance, created_at, updated_at 
		FROM users
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Number, &u.Address, &u.Balance, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

func (r *repository) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	sqlStatement :=
		`SELECT id, name, price, sku, amount
		FROM products
		WHERE id = $1
`
	product := &domain.Product{}

	err := r.conn.QueryRow(ctx, sqlStatement, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.SKU,
		&product.Amount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrorProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (r *repository) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	sqlStatement := `
		SELECT id, name, price, sku, amount
		FROM products
		ORDER BY id ASC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.conn.Query(ctx, sqlStatement, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.SKU, &p.Amount, &p.CategoryID)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return []*domain.Product{}, nil
	}
	return products, nil
}

func (r *repository) CreateOrder(ctx context.Context, order *domain.Order) error {
	sqlStatement := `
		INSERT INTO orders (id, user_id, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
`
	err := r.conn.QueryRow(
		ctx,
		sqlStatement,
		order.ID,
		order.UserID,
		order.CreatedAt,
		order.UpdatedAt,
		order.Status,
	).Scan(&order.ID)

	return err
}

func (r *repository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	sqlStatement := `
		UPDATE orders
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id;
`
	err := r.conn.QueryRow(
		ctx,
		sqlStatement,
		order.Status,
		order.ID,
	).Scan(&order.ID)
	return err

}

func (r *repository) GetOrderByID(ctx context.Context, ID string) (*domain.Order, error) {
	sqlStatement := `
SELECT id, userid, created_at, updated_at, status
FROM orders
WHERE id = $1`
	order := &domain.Order{}
	err := r.conn.QueryRow(ctx, sqlStatement, ID).Scan(
		&order.ID,
		&order.UserID,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrorProductNotFound
		}
		return nil, err
	}

	return nil, err
}

func (r *repository) AddProductToCategory(ctx context.Context, categoryID, productID string) error {
	sqlStatemnt := `
		UPDATE products
		SET category_id = $1
		WHERE id = $2
		RETURNING id;
	`

	var id string
	err := r.conn.QueryRow(
		ctx,
		sqlStatemnt,
		categoryID,
		productID,
	).Scan(&id)
	return err
}

func (r *repository) RemoveProductFromCategory(ctx context.Context, categoryID, productID string) error {
	sqlStatement := `
		UPDATE products
		SET product_id = NULL
		WHERE id = $2
		RETURNING id;
	`
	var id string
	err := r.conn.QueryRow(
		ctx,
		sqlStatement,
		categoryID,
		productID,
	).Scan(&id)
	return err
}

func (r *repository) GetSupplierByID(ctx context.Context, ID string) (*domain.Supplier, error) {
	sqlStatement := `
		SELECT id, name
		FROM suppliers
		WHERE id = $1
`
	supplier := &domain.Supplier{}
	err := r.conn.QueryRow(ctx, sqlStatement, ID).Scan(&supplier.ID, &supplier.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrorSupplierNotFound
		}
		return nil, err
	}
	return supplier, nil
}

func (r *repository) DeleteSupplierByID(ctx context.Context, ID string) error {
	sqlStatement := `
	DELETE FROM suppliers
	WHERE id = $1
	`
	tag, err := r.conn.Exec(ctx, sqlStatement, ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrorSupplierNotFound
	}
	return nil
}
