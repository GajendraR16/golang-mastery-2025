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
	"google.golang.org/grpc/status"
)

type taskServer struct {
	pb.UnimplementedTaskServiceServer
	mu     sync.Mutex
	tasks  map[int32]*pb.Task
	nextID int32
}

func newTaskServer() *taskServer {
	return &taskServer{
		tasks:  make(map[int32]*pb.Task),
		nextID: 1,
	}
}

func loggingUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resq, err := handler(ctx, req)
	duration := time.Since(start)
	if err != nil {
		log.Printf("RPC Failed | Method: %s | Duration: %v | Error: %v", info.FullMethod, duration, err)
	} else {
		log.Printf("RPC Success | Method: %s | Duration: %v", info.FullMethod, duration)
	}
	return resq, err
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description cannot be empty")
	}

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

	log.Printf("Created task: id=%d description=%s", task.Id, task.Description)

	return &pb.CreateTaskResponse{Task: task}, nil
}

func (s *taskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "task %d not found", req.Id)
	}

	return &pb.GetTaskResponse{Task: task}, nil
}

func (s *taskServer) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks := make([]*pb.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return &pb.ListTasksResponse{Tasks: tasks}, nil
}

func (s *taskServer) CompleteTasks(ctx context.Context, req *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "task %d not found", req.Id)
	}

	task.Completed = true
	task.CompletedAt = time.Now().Format(time.RFC3339)

	return &pb.CompleteTaskResponse{Task: task}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(loggingUnaryInterceptor))
	pb.RegisterTaskServiceServer(grpcServer, newTaskServer())

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
