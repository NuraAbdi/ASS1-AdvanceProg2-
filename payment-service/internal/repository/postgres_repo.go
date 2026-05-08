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

func (r *PostgresRepo) GetStats() (int64, int64, int64, int64, error) {
	query := `
	SELECT 
		COUNT(*) as total,
		SUM(CASE WHEN status='Authorized' THEN 1 ELSE 0 END),
		SUM(CASE WHEN status='Declined' THEN 1 ELSE 0 END),
		COALESCE(SUM(amount), 0)
	FROM payments
	`

	var total, authorized, declined, amount int64

	err := r.db.QueryRow(query).Scan(&total, &authorized, &declined, &amount)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return total, authorized, declined, amount, nil
}
