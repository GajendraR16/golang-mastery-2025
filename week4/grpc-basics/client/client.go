package main

import (
	"context"
	"log"
	"time"

	pb "grpc-basics/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	//Create Task
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	createResp, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
		Description: "Learn gRPC",
	})

	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	log.Printf("Created task: %v", createResp.Task)

	//List Task
	listResp, err := client.ListTask(ctx, &pb.ListTaskRequest{})

	if err != nil {
		log.Fatalf("could not list tasks: %v", err)
	}
	log.Printf("Tasks: %v", listResp.Tasks)
}
