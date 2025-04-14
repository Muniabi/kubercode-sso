FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o kubercode-sso ./cmd/sso/main.go

FROM alpine:latest

LABEL org.opencontainers.image.title="kubercode-sso"
LABEL org.opencontainers.image.description="SSO Service for Kubercode"
LABEL org.opencontainers.image.vendor="Kubercode"

WORKDIR /root/

COPY --from=builder /app/kubercode-sso .
COPY --from=builder /app/certs ./certs/
COPY --from=builder /app/config/config.yaml ./config/

EXPOSE 50051
EXPOSE 1488

CMD ["./kubercode-sso"]
