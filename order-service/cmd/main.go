package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	pb "order-service/proto/paymentpb"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"order-service/internal/cache"
	"order-service/internal/paymentclient"
	"order-service/internal/repository"
	grpcTransport "order-service/internal/transport/grpc"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=2007 dbname=order_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresRepo(db)

	redisCache := cache.NewRedisCache("localhost:6379")

	paymentAddr := os.Getenv("GRPC_PAYMENT_ADDR")

	paymentClient, err := paymentclient.NewGRPCClient(paymentAddr)
	if err != nil {
		log.Fatal(err)
	}

	uc := usecase.NewOrderUsecase(repo, paymentClient, redisCache)

	// 🔹 HTTP
	go func() {
		r := gin.Default()
		handler := http.NewHandler(uc)

		r.POST("/orders", handler.CreateOrder)
		r.GET("/orders/:id", handler.GetOrder)
		r.PATCH("/orders/:id/cancel", handler.CancelOrder)
		r.GET("/orders/stats", handler.GetStats)

		log.Println("Order HTTP running on :8080")

		if err := r.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	// 🔹 gRPC (streaming)
	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatal(err)
		}

		server := grpcTransport.NewServer(uc)

		grpcServer := grpc.NewServer()
		pb.RegisterOrderServiceServer(grpcServer, server)

		log.Println("Order gRPC running on :50052")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
