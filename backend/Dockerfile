FROM golang:1.24.1-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify && go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

RUN go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["sh", "-c", "./server"]
