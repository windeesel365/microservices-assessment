
FROM golang:1.22.3 AS builder

WORKDIR /app

COPY ./services/user/go.mod ./services/user/go.sum ./

RUN go mod download


COPY ../../proto /app/proto
COPY . .

COPY ./services/user/.env .env

RUN CGO_ENABLED=0 GOOS=linux go build -v -o myapp ./services/user/server.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./myapp"]
