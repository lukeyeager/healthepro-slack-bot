FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:3.21
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
COPY --from=builder /build/bot .
ENTRYPOINT ["./bot"]
