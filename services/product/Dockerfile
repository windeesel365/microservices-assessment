
FROM golang:1.22.3 AS builder

WORKDIR /app

COPY ./services/product/go.mod ./services/product/go.sum ./

RUN go mod download
RUN go mod tidy

COPY ../../proto /app/proto
COPY . .

COPY ./services/product /app/services/product

COPY ./services/product/.env .env

RUN CGO_ENABLED=0 GOOS=linux go build -v -o myapp2 ./services/product/server.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp2 .
COPY --from=builder /app/.env .

EXPOSE 8081

CMD ["./myapp2"]
