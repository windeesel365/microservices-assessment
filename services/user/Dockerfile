FROM golang:1.22.3 AS builder

WORKDIR /app/services/user

COPY ./services/user/go.mod .
COPY ./services/user/go.sum .

RUN go mod download

COPY ./services/user/ .
COPY ../../proto /app/proto


RUN go get github.com/jmoiron/sqlx
RUN go get github.com/joho/godotenv
RUN go get github.com/lib/pq
RUN go get golang.org/x/crypto/bcrypt
RUN go get google.golang.org/grpc

RUN go mod tidy


RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/myapp1 ./server.go



FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp1 .
COPY --from=builder /app/services/user/.env .

EXPOSE 8080

CMD ["./myapp1"]
