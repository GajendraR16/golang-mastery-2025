# 1. Install Protocol Buffer compiler
# macOS:
brew install protobuf

# Ubuntu/Linux:
sudo apt install protobuf-compiler

# Windows:
# Download from https://github.com/protocolbuffers/protobuf/releases

# 2. Verify installation
protoc --version
# Should show: libprotoc 3.x or higher

# 3. Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 4. Update PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$PATH:$(go env GOPATH)/bin"

# 5. Verify plugins
which protoc-gen-go
which protoc-gen-go-grpc
# Both should show paths

# 6. Install Redis (for Day 4-5)
# macOS:
brew install redis

# Ubuntu:
sudo apt install redis-server

# 7. Test Redis
redis-cli ping
# Should respond: PONG

# Code generation proto command
protoc --go_out=. --go_opt=paths=source_relative\
--go-grpc_out=. --go-grpc_opt=paths=source_relative\
proto/task.proto

# proto/task.proto is the proto file / . to generate in the current directory

# Using evans

# 8. Install evans (gRPC client)
# macOS:
brew install evans

# Linux:
# Download evans package using curl:

curl -L https://github.com/ktr0731/evans/releases/latest/download/evans_linux_amd64.tar.gz -o evans.tar.gz
tar -xzf evans.tar.gz
chmod +x evans

# Verify installation
./evans --version

# 9. Start the gRPC server
cd week4/grpc-basics
go run ./server

# 10. Connect with evans using server reflection
./evans --host localhost --port 50051 repl

# Useful evans commands inside the REPL
show package
package <task>
show service
service TaskService
show message

# 11. Call RPCs from evans
# CreateTask
call CreateTask
# Enter:
# {
#   "description": "Learn gRPC with evans"
# }

# GetTask
call GetTask
# Enter:
# {
#   "id": 1
# }

# ListTask
call ListTask
# Enter:
# {}

# 12. If reflection is disabled, use the proto file directly
./evans --host localhost --port 50051 -r repl \
  --proto proto/task.proto
