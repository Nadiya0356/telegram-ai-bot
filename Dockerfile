FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bot ./cmd/bot

FROM alpine:latest
COPY --from=builder /app/bot /bot
CMD ["/bot"]