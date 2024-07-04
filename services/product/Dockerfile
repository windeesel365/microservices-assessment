FROM golang:1.22.3 AS builder

WORKDIR /app/services/product

COPY ./services/user/go.mod .
COPY ./services/user/go.sum .

RUN go mod download

COPY ./services/product/ .
COPY ../../proto /app/proto

RUN go get -v ./...

RUN go get github.com/jmoiron/sqlx
RUN go get github.com/joho/godotenv
RUN go get github.com/lib/pq
RUN go get golang.org/x/crypto/bcrypt
RUN go get google.golang.org/grpc

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/myapp2 ./server.go



FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp2 .
COPY --from=builder /app/services/product/.env .

EXPOSE 8081

CMD ["./myapp2"]
