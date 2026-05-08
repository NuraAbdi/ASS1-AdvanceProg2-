package paymentclient

import (
	"context"
	"time"

	pb "order-service/proto/paymentpb"

	"google.golang.org/grpc"
)

type GRPCClient struct {
	client pb.PaymentServiceClient
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := pb.NewPaymentServiceClient(conn)

	return &GRPCClient{client: c}, nil
}

func (c *GRPCClient) Pay(orderID string, amount int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := c.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}
