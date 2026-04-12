FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(git describe --tags --always 2>/dev/null || echo dev)" -o /dockguard ./cmd/dockguard/

FROM alpine:3.20
COPY --from=builder /dockguard /usr/local/bin/dockguard
ENTRYPOINT ["dockguard"]
