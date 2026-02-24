FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o redirect-service .

FROM alpine:3.19

COPY --from=builder /app/redirect-service /redirect-service

ENTRYPOINT ["/redirect-service"]
