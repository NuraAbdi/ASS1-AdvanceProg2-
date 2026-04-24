package main

import (
	"log"
	"net"

	"payment-service/internal/domain"
	grpcTransport "payment-service/internal/transport/grpc"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"
	pb "payment-service/proto/paymentpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type InMemoryRepo struct{}

func (r *InMemoryRepo) GetStats() (int64, int64, int64, int64, error) {
	return 0, 0, 0, 0, nil
}
func (r *InMemoryRepo) Save(p *domain.Payment) error {
	return nil
}

func main() {
	repo := &InMemoryRepo{}
	uc := usecase.NewPaymentUsecase(repo)

	// 🔹 HTTP server (REST)
	go func() {
		r := gin.Default()
		handler := http.NewHandler(uc)

		r.POST("/payments", handler.CreatePayment)

		log.Println("HTTP running on :8081")

		if err := r.Run(":8081"); err != nil {
			log.Fatal(err)
		}
	}()

	// 🔹 gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal(err)
		}

		server := grpcTransport.NewServer(uc)

		grpcServer := grpc.NewServer()
		pb.RegisterPaymentServiceServer(grpcServer, server)

		log.Println("gRPC running on :50051")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	// 🔒 чтобы программа не завершалась
	select {}
}
