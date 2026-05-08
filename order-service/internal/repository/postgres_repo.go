package repository

import (
	"database/sql"
	"order-service/internal/domain"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(o *domain.Order) error {
	query := `
	INSERT INTO orders (id, customer_id, item_name, amount, status, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query,
		o.ID,
		o.CustomerID,
		o.ItemName,
		o.Amount,
		o.Status,
		o.CreatedAt,
	)

	return err
}

func (r *PostgresRepo) Update(o *domain.Order) error {
	query := `UPDATE orders SET status=$1 WHERE id=$2`
	_, err := r.db.Exec(query, o.Status, o.ID)
	return err
}

func (r *PostgresRepo) GetByID(id string) (*domain.Order, error) {
	query := `
	SELECT id, customer_id, item_name, amount, status, created_at
	FROM orders WHERE id=$1
	`

	row := r.db.QueryRow(query, id)

	order := &domain.Order{}

	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.ItemName,
		&order.Amount,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *PostgresRepo) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	return count, err
}

func (r *PostgresRepo) CountByStatus(status string) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE status=$1", status).Scan(&count)
	return count, err
}
