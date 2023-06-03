env:
	go mod download -x

player:
	go run cmd/client/main.go $(NAME)

bot:
	go run cmd/client/main.go auto