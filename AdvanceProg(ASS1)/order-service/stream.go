package main

import (
	"context"
	"log"

	pb "order-service/proto/paymentpb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewOrderServiceClient(conn)

	stream, err := client.SubscribeToOrderUpdates(context.Background(), &pb.OrderRequest{
		OrderId: "97f44b6b-ba10-401a-ad47-4713672b4ba2",
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("STATUS:", res.Status)
	}
}
