package main

import (
	"context"
	"log"
	"net"
	"strings"
	"sync"

	pb "user-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServer struct {
	pb.UnimplementedUserServiceServer

	mu     sync.Mutex
	users  map[int32]*pb.User
	nextID int32
}

func (s *userServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		return nil, status.Error(codes.InvalidArgument, "name and email are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user := &pb.User{
		Id:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
	}

	s.users[s.nextID] = user
	s.nextID++

	return &pb.CreateUserResponse{
		User: user,
	}, nil
}

func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
	}

	return &pb.GetUserResponse{
		User: user,
	}, nil
}

func (s *userServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]*pb.User, 0, len(s.users))

	for _, user := range s.users {
		users = append(users, user)
	}

	return &pb.ListUsersResponse{
		Users: users,
	}, nil
}

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("method=%s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("method=%s error=%v", info.FullMethod, err)
	}

	return resp, err
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	pb.RegisterUserServiceServer(server, &userServer{
		users:  make(map[int32]*pb.User),
		nextID: 1,
	})

	log.Println("gRPC server listening on :50052")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
