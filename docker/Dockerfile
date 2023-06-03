FROM golang:1.18-alpine

WORKDIR /server
COPY go.* /server/.

RUN go mod download -x


ENTRYPOINT ["go", "run", "cmd/server/main.go"]