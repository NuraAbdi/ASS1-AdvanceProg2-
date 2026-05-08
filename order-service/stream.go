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
		OrderId: "4573976c-16fd-4a73-891f-3c88edbdf60f",
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
