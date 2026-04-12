package grpc

import (
	"time"

	"order-service/internal/usecase"
	pb "order-service/proto/paymentpb"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	uc *usecase.OrderUsecase
}

func NewServer(uc *usecase.OrderUsecase) *Server {
	return &Server{uc: uc}
}

func (s *Server) SubscribeToOrderUpdates(
	req *pb.OrderRequest,
	stream pb.OrderService_SubscribeToOrderUpdatesServer,
) error {

	orderID := req.OrderId
	var lastStatus string

	for {
		order, err := s.uc.GetOrder(orderID)
		if err != nil {
			return err
		}

		// если статус изменился → отправляем
		if order.Status != lastStatus {
			lastStatus = order.Status

			stream.Send(&pb.OrderStatusUpdate{
				OrderId: order.ID,
				Status:  order.Status,
			})
		}

		time.Sleep(1 * time.Second)
	}
}
