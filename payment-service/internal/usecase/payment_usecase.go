package usecase

import (
	"payment-service/internal/domain"
	"payment-service/internal/messaging"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Save(payment *domain.Payment) error
	GetStats() (int64, int64, int64, int64, error)
}

type PaymentUsecase struct {
	repo      PaymentRepository
	publisher *messaging.Publisher
}

func NewPaymentUsecase(r PaymentRepository, p *messaging.Publisher) *PaymentUsecase {
	return &PaymentUsecase{
		repo:      r,
		publisher: p,
	}
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

	if payment.Status == "Authorized" {
		uc.publisher.Publish(payment.OrderID, payment.Amount)
	}
	
	return payment, nil
}

func (uc *PaymentUsecase) GetStats() (int64, int64, int64, int64, error) {
	return uc.repo.GetStats()
}
