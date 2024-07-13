FROM golang:1.22.3 AS builder

WORKDIR /app/services/payment

COPY ./services/payment/go.mod .
COPY ./services/payment/go.sum .

RUN go mod download

COPY ./services/payment/ .
COPY ../../proto /app/proto

RUN go get -v ./...

RUN go get github.com/jmoiron/sqlx
RUN go get github.com/joho/godotenv
RUN go get github.com/lib/pq
RUN go get google.golang.org/grpc
RUN go get google.golang.org/grpc/credentials/insecure

RUN go get github.com/windeesel365/microservices-assessment/services/order


RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/myapp4 ./server.go



FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp4 .
COPY --from=builder /app/services/payment/.env .

EXPOSE 8083

CMD ["./myapp4"]
