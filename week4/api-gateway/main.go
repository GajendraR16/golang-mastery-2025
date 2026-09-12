package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	pb "grpc-basics/proto"
	ub "user-service/proto"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Gateway struct {
	taskClient pb.TaskServiceClient
	userClient ub.UserServiceClient
}

type TaskData struct {
	Description string `json:"description"`
	UserId      int    `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TaskWithUser struct {
	Task *pb.Task `json:"task"`
	User *ub.User `json:"user"`
}

type TasksResponse struct {
	Tasks []*TaskWithUser `json:"tasks"`
}

var ErrorCodeMapping = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusBadRequest,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unauthenticated:    http.StatusUnauthorized,
}

func main() {

	taskaddr := os.Getenv("TASK_SERVICE_ADDR")
	if taskaddr == "" {
		taskaddr = "localhost:50051" // default for local dev
	}

	taskconn, err := grpc.NewClient(
		taskaddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer taskconn.Close()

	useraddr := os.Getenv("USER_SERVICE_ADDR")
	if useraddr == "" {
		useraddr = "localhost:50052" // default for local dev
	}

	userconn, err := grpc.NewClient(
		useraddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer userconn.Close()

	//Generate gRPC client
	taskclient := pb.NewTaskServiceClient(taskconn)
	userclient := ub.NewUserServiceClient(userconn)

	gateway := &Gateway{
		taskClient: taskclient,
		userClient: userclient,
	}

	//rest router
	r := mux.NewRouter()

	//Tasks
	r.HandleFunc("/tasks", gateway.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", gateway.GetTaskById).Methods("GET")
	r.HandleFunc("/tasks", gateway.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", gateway.CompleteTaskById).Methods("PUT")

	//Users
	r.HandleFunc("/users", gateway.CreateUsers).Methods("POST")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func jsonError(w http.ResponseWriter, message string, code int, err error) {
	if err != nil && code == 500 {
		slog.Error(message, "error", err.Error())
	} else {
		slog.Warn(message, "status", code)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func jsonHandler(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to Encode Json", http.StatusInternalServerError)
		return
	}
}

func writeProto(w http.ResponseWriter, code int, msg proto.Message) {
	b, err := protojson.Marshal(msg)
	if err != nil {
		jsonError(w, "failed to encode response", http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(b)

}

func (g *Gateway) GetUser(w http.ResponseWriter, r *http.Request, id int32) (*ub.User, error) {
	userReq := &ub.GetUserRequest{Id: id}
	userResp, err := g.userClient.GetUser(r.Context(), userReq)
	if err != nil {
		return nil, err
	}
	return userResp.User, nil
}

func (g *Gateway) GetTask(w http.ResponseWriter, r *http.Request) {
	req := &pb.ListTasksRequest{}

	resp, err := g.taskClient.ListTasks(r.Context(), req)
	if err != nil {
		code := status.Code(err)
		httpCode, ok := ErrorCodeMapping[code]
		if !ok {
			httpCode = http.StatusInternalServerError
		}

		jsonError(w, err.Error(), httpCode, err)
		return
	}

	// userID -> User
	users := make(map[int32]*ub.User)

	var taskresponse []*TaskWithUser

	for _, task := range resp.Tasks {
		user, ok := users[task.UserId]

		if !ok {
			user, err = g.GetUser(w, r, task.UserId)
			if err != nil {
				code := status.Code(err)
				httpCode := ErrorCodeMapping[code]
				jsonError(w, err.Error(), httpCode, err)
				return
			}

			users[task.UserId] = user
		}

		taskWithUser := &TaskWithUser{
			Task: task,
			User: user,
		}

		taskresponse = append(taskresponse, taskWithUser)
	}

	jsonHandler(w, http.StatusOK, &TasksResponse{
		Tasks: taskresponse,
	})
}

func (g *Gateway) GetTaskById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, "Invalid ID", http.StatusBadRequest, err)
		return
	}

	req := &pb.GetTaskRequest{
		Id: int32(id),
	}

	resp, err := g.taskClient.GetTask(r.Context(), req)
	if err != nil {
		code := status.Code(err)
		httpCode, ok := ErrorCodeMapping[code]

		if !ok {
			httpCode = http.StatusInternalServerError
		}

		jsonError(w, err.Error(), httpCode, err)
		return
	}

	user, err := g.GetUser(w, r, resp.Task.UserId)
	if err != nil {
		code := status.Code(err)
		httpCode := ErrorCodeMapping[code]
		jsonError(w, err.Error(), httpCode, err)
		return
	}

	taskWithUser := &TaskWithUser{
		Task: resp.Task,
		User: user,
	}

	jsonHandler(w, http.StatusOK, taskWithUser)
}

func (g *Gateway) CompleteTaskById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, "Invalid ID", http.StatusBadRequest, err)
		return
	}

	req := &pb.CompleteTaskRequest{Id: int32(id)}

	resp, err := g.taskClient.CompleteTask(r.Context(), req)
	if err != nil {
		code := status.Code(err)
		httpCode, ok := ErrorCodeMapping[code]

		if !ok {
			httpCode = http.StatusInternalServerError
		}

		jsonError(w, err.Error(), httpCode, err)
		return
	}

	writeProto(w, http.StatusOK, resp)
}

func (g *Gateway) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task TaskData
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		jsonError(w, "Invalid Json", http.StatusBadRequest, err)
		return
	}
	_, err := g.GetUser(w, r, int32(task.UserId))
	if err != nil {
		code := status.Code(err)
		httpCode := ErrorCodeMapping[code]
		jsonError(w, err.Error(), httpCode, err)
		return
	}
	req := &pb.CreateTaskRequest{Description: task.Description, UserId: int32(task.UserId)}

	resp, err := g.taskClient.CreateTask(r.Context(), req)
	if err != nil {
		code := status.Code(err)
		httpCode, ok := ErrorCodeMapping[code]

		if !ok {
			httpCode = http.StatusInternalServerError
		}

		jsonError(w, err.Error(), httpCode, err)
		return
	}

	writeProto(w, http.StatusCreated, resp)
}

func (g *Gateway) CreateUsers(w http.ResponseWriter, r *http.Request) {
	var user UserData
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jsonError(w, "Invalid Json", http.StatusBadRequest, err)
		return
	}

	req := &ub.CreateUserRequest{Name: user.Name, Email: user.Email}
	resp, err := g.userClient.CreateUser(r.Context(), req)
	if err != nil {
		code := status.Code(err)
		httpCode, ok := ErrorCodeMapping[code]

		if !ok {
			httpCode = http.StatusInternalServerError
		}

		jsonError(w, err.Error(), httpCode, err)
		return
	}
	writeProto(w, http.StatusCreated, resp)
}
