FROM golang:1.23-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bbb main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/bbb .
