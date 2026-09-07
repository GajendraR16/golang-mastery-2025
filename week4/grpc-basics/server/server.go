package main

import (
	"context"
	pb "grpc-basics/proto"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type taskServer struct {
	pb.UnimplementedTaskServiceServer
	mu     sync.Mutex
	tasks  map[int32]*pb.Task
	nextID int32
}

func NewTaskServer() *taskServer {
	return &taskServer{
		tasks:  make(map[int32]*pb.Task),
		nextID: 1,
	}
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &pb.Task{
		Id:          s.nextID,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	s.tasks[s.nextID] = task
	s.nextID++

	log.Printf("Created task: %v", task)
	return &pb.CreateTaskResponse{Task: task}, nil
}

func (s *taskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "task not found")
	}

	return &pb.GetTaskResponse{Task: task}, nil
}

func (s *taskServer) ListTask(ctx context.Context, req *pb.ListTaskRequest) (*pb.ListTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks := make([]*pb.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return &pb.ListTaskResponse{Tasks: tasks}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Initialize gRPC Engine
	grpcServer := grpc.NewServer()

	//Register's Task Service
	pb.RegisterTaskServiceServer(grpcServer, NewTaskServer())

	//Register's all the method's which helps in calling using cmd.
	//Instead of hard coding in client
	//For more options read README.md
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
