FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/health-monitor ./cmd/healer

FROM alpine:3.22

RUN apk --no-cache add ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /out/health-monitor /usr/local/bin/health-monitor

USER app
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/health-monitor"]
