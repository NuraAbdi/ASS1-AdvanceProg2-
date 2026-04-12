package usecase

import (
	"fmt"
	"time"

	"order-service/internal/domain"

	"github.com/google/uuid"
)

type PaymentClient interface {
	Pay(orderID string, amount int64) (string, error)
}

type OrderRepository interface {
	Save(order *domain.Order) error
	Update(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	CountAll() (int, error)
	CountByStatus(status string) (int, error)
}

type OrderUsecase struct {
	repo    OrderRepository
	payment PaymentClient
}

func NewOrderUsecase(r OrderRepository, p PaymentClient) *OrderUsecase {
	return &OrderUsecase{
		repo:    r,
		payment: p,
	}
}

func (uc *OrderUsecase) CreateOrder(customerID, itemName string, amount int64) (*domain.Order, error) {
	order := &domain.Order{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     "Pending",
		CreatedAt:  time.Now(),
	}

	// сохраняем заказ
	err := uc.repo.Save(order)
	if err != nil {
		return nil, err
	}

	// вызываем Payment Service через интерфейс (правильно!)
	//status, err := uc.payment.Pay(order.ID, order.Amount)
	//if err != nil {
	//	order.Status = "Failed"
	//	_ = uc.repo.Update(order)
	//	return nil, err
	//}
	//
	//if status == "Authorized" {
	//	order.Status = "Paid"
	//} else {
	//	order.Status = "Failed"
	//}

	order.Status = "Pending"
	_ = uc.repo.Update(order)

	// обновляем статус
	_ = uc.repo.Update(order)

	return order, nil
}

func (uc *OrderUsecase) GetOrder(id string) (*domain.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *OrderUsecase) CancelOrder(id string) error {
	order, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status != "Pending" {
		return fmt.Errorf("only pending orders can be cancelled")
	}

	order.Status = "Cancelled"
	return uc.repo.Update(order)
}

func (uc *OrderUsecase) GetStats() (map[string]int, error) {
	total, _ := uc.repo.CountAll()
	pending, _ := uc.repo.CountByStatus("Pending")
	paid, _ := uc.repo.CountByStatus("Paid")
	failed, _ := uc.repo.CountByStatus("Failed")
	cancelled, _ := uc.repo.CountByStatus("Cancelled")

	return map[string]int{
		"total":     total,
		"pending":   pending,
		"paid":      paid,
		"failed":    failed,
		"cancelled": cancelled,
	}, nil
}

func (uc *OrderUsecase) GetStatsByStatus(status string) (int, error) {
	valid := map[string]bool{
		"pending":   true,
		"paid":      true,
		"failed":    true,
		"cancelled": true,
	}

	if !valid[status] {
		return -1, fmt.Errorf("invalid status")
	}

	// capitalize
	statusMap := map[string]string{
		"pending":   "Pending",
		"paid":      "Paid",
		"failed":    "Failed",
		"cancelled": "Cancelled",
	}

	count, err := uc.repo.CountByStatus(statusMap[status])
	if err != nil {
		return 0, err
	}

	return count, nil
}
