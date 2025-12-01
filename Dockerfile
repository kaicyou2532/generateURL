# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /server ./server
RUN mkdir -p /app/uploads

ENV PORT=8080 \
	UPLOAD_DIR=/app/uploads

EXPOSE 8080
ENTRYPOINT ["/app/server"]
