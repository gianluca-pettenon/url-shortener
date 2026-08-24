FROM golang:1.27.0 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o shortener ./cmd/shortener

FROM alpine:3.24.1
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/shortener .
EXPOSE 80

CMD ["./shortener", "serve"]