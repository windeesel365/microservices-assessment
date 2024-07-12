FROM golang:1.22.3 AS builder

WORKDIR /app/services/order

COPY ./services/order/go.mod .
COPY ./services/order/go.sum .

RUN go mod download

COPY ./services/order/ .
COPY ../../proto /app/proto

RUN go get -v ./...

RUN go get github.com/jmoiron/sqlx
RUN go get github.com/joho/godotenv
RUN go get github.com/lib/pq
RUN go get google.golang.org/grpc
RUN go get google.golang.org/grpc/credentials/insecure

RUN go get github.com/windeesel365/microservices-assessment/services/user
RUN go get github.com/windeesel365/microservices-assessment/services/product

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/myapp3 ./server.go



FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp3 .
COPY --from=builder /app/services/order/.env .

EXPOSE 8082

CMD ["./myapp3"]
