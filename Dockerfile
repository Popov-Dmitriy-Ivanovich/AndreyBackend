FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
ENV CGO_ENABLED=1
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

RUN go build -o run main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/run /build/run
EXPOSE 8080
CMD ["./run"]