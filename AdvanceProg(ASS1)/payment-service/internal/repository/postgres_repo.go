package repository

import (
	"database/sql"
	"payment-service/internal/domain"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(p *domain.Payment) error {
	query := `
	INSERT INTO payments (id, order_id, transaction_id, amount, status)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query,
		p.ID,
		p.OrderID,
		p.TransactionID,
		p.Amount,
		p.Status,
	)

	return err
}
