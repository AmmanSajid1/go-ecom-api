FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .
COPY internal/adapters/postgresql/migrations ./internal/adapters/postgresql/migrations

EXPOSE 8080

CMD ["sh", "-c", "./api migrate && ./api"]

