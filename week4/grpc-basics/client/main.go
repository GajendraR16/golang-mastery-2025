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
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Create task
	createResp, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
		Description: "Learn gRPC"})
	if err != nil {
		log.Fatalf("CreateTask Failed: %v", err)
	}
	log.Printf("Created %v", createResp.Task)

	//Get Task
	getResp, err := client.GetTask(ctx, &pb.GetTaskRequest{Id: createResp.Task.Id})
	if err != nil {
		log.Fatalf("GetTask Failed: %v", err)
	}
	log.Printf("Got: %v", getResp.Task)

	//List Task
	listResp, err := client.ListTasks(ctx, &pb.ListTasksRequest{})
	if err != nil {
		log.Fatalf("ListTask Failed: %v", err)
	}
	log.Printf("All tasks: %v", listResp.Tasks)

	//Complete Task
	completeResp, err := client.CompleteTask(ctx, &pb.CompleteTaskRequest{Id: createResp.Task.Id})
	if err != nil {
		log.Fatalf("CompleteTask Failed: %v", err)
	}
	log.Printf("Complete task: %v", completeResp)

	//Test error handling - get non-existent task
	_, err = client.GetTask(ctx, &pb.GetTaskRequest{Id: 920})
	if err != nil {
		log.Printf("Expected error: %v", err)
	}
}
