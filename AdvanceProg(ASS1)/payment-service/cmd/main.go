package main

import (
	"payment-service/internal/domain"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type InMemoryRepo struct{}

func (r *InMemoryRepo) Save(p *domain.Payment) error {
	return nil
}

func main() {
	r := gin.Default()

	repo := &InMemoryRepo{}
	uc := usecase.NewPaymentUsecase(repo)
	handler := http.NewHandler(uc)

	r.POST("/payments", handler.CreatePayment)

	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
