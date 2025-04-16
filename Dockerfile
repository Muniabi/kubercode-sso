FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o kubercode-sso ./cmd/main.go

FROM debian:bookworm-slim

LABEL org.opencontainers.image.title="kubercode-sso"
LABEL org.opencontainers.image.description="SSO Service for Kubercode"
LABEL org.opencontainers.image.vendor="Kubercode"

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/kubercode-sso .
COPY --from=builder /app/certs ./certs/
COPY --from=builder /app/config/config.yaml ./config/

EXPOSE 50051
EXPOSE 1488

CMD ["./kubercode-sso"]
