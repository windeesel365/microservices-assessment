
FROM golang:1.22.3 AS builder

WORKDIR /app

COPY ./services/user/go.mod ./services/user/go.sum ./

RUN go mod download
RUN go mod tidy

COPY ../../proto /app/proto
COPY . .
COPY ./services/user /app/services/user
COPY ./services/user/.env .env

RUN CGO_ENABLED=0 GOOS=linux go build -v -o user_service ./services/user/server.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/user_service .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./user_service"]
