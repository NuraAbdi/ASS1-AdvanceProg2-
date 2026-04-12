package main

import (
	"database/sql"
	"log"

	"order-service/internal/paymentclient"
	"order-service/internal/repository"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()

	connStr := "host=localhost port=5432 user=postgres password=2007 dbname=order_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresRepo(db)
	paymentClient := paymentclient.NewClient()

	uc := usecase.NewOrderUsecase(repo, paymentClient)
	handler := http.NewHandler(uc)

	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders/:id", handler.GetOrder)
	r.PATCH("/orders/:id/cancel", handler.CancelOrder)
	r.GET("/orders/stats", handler.GetStats)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
