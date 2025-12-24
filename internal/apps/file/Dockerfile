# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

LABEL stage=builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o kgym-file-service ./cmd/main.go

FROM alpine:latest

LABEL maintainer="Alexandr Rutkowski <kitanoyoru@icloud.com>"
WORKDIR /usr/local/bin

COPY --from=builder /build/kgym-file-service .

ENTRYPOINT ["./kgym-file-service"]

CMD ["run"]
