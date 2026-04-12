package usecase

import (
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Save(payment *domain.Payment) error
}

type PaymentUsecase struct {
	repo PaymentRepository
}

func NewPaymentUsecase(r PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: r}
}

func (uc *PaymentUsecase) ProcessPayment(orderID string, amount int64) (*domain.Payment, error) {
	payment := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
	}

	// business rule
	if amount > 100000 {
		payment.Status = "Declined"
	} else {
		payment.Status = "Authorized"
	}

	err := uc.repo.Save(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
