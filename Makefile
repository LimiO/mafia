dirs:
	mkdir content
	mkdir content/pdf
	mkdir content/img

env:
	go mod download -x
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/game/actions.proto && \
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/game/roles.proto && \
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/game/state.proto && \
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/connection/connection.proto

player:
	go run cmd/client/main.go

bot:
	go run cmd/client/main.go auto