package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"order-service/internal/cache"
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
	cache   *cache.RedisCache
}

func NewOrderUsecase(r OrderRepository, p PaymentClient, c *cache.RedisCache) *OrderUsecase {
	return &OrderUsecase{
		repo:    r,
		payment: p,
		cache:   c,
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
	status, err := uc.payment.Pay(order.ID, order.Amount)
	if err != nil {
		order.Status = "Failed"
		_ = uc.repo.Update(order)

		_ = uc.cache.Delete("order:" + order.ID)

		return nil, err
	}

	if status == "Authorized" {
		order.Status = "Paid"
	} else {
		order.Status = "Failed"
	}

	// обновляем статус
	_ = uc.repo.Update(order)

	_ = uc.cache.Delete("order:" + order.ID)

	return order, nil
}

func (uc *OrderUsecase) GetOrder(id string) (*domain.Order, error) {

	// 1. check redis cache
	cached, err := uc.cache.Get("order:" + id)

	if err == nil {
		log.Println("CACHE HIT")

		var order domain.Order

		err = json.Unmarshal([]byte(cached), &order)
		if err == nil {
			return &order, nil
		}
	}

	log.Println("CACHE MISS")

	// 2. get from postgres
	order, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 3. save to redis
	data, _ := json.Marshal(order)

	_ = uc.cache.Set(
		"order:"+id,
		string(data),
	)

	return order, nil
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
	err = uc.repo.Update(order)
	if err != nil {
		return err
	}

	// invalidate redis cache
	_ = uc.cache.Delete("order:" + id)

	return nil
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
