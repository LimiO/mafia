env:
	go mod download -x

player:
	go run cmd/client/main.go

bot:
	go run cmd/client/main.go auto