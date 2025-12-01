# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server
RUN mkdir -p /uploads

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /server ./server
COPY --from=builder --chown=nonroot:nonroot /uploads ./uploads

ENV PORT=8000 \
	UPLOAD_DIR=/app/uploads

EXPOSE 8000
ENTRYPOINT ["/app/server"]
