
FROM golang:1.22.3 AS builder

WORKDIR /app

COPY ./services/product/go.mod ./services/product/go.sum ./

RUN go mod download


COPY ../../proto /app/proto
COPY . .

COPY ./services/product/.env .env

RUN CGO_ENABLED=0 GOOS=linux go build -v -o myapp ./services/product/server.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

EXPOSE 8081

CMD ["./myapp"]
