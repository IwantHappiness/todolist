FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app

FROM alpine:3.22

WORKDIR /app

RUN adduser -D appuser

COPY --from=builder /app/app .

USER appuser

EXPOSE 8080

CMD ["./app"]
