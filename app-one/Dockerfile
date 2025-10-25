FROM golang:1.25-alpine AS builder
RUN apk --no-cache add git ca-certificates tzdata

ARG APP_VERSION=dev
ARG GO_VERSION=1.25.2


WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${APP_VERSION}" \
    -o /build/server \
    ./cmd/main.go


FROM alpine:latest
RUN apk --no-cache add git ca-certificates tzdata


WORKDIR /app

COPY --from=builder /build/server .

RUN adduser -D -s /bin/bash appuser
USER appuser

EXPOSE 8081
ENTRYPOINT ["./server"]