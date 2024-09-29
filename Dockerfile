FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o run main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/run /build/run
EXPOSE 8080
CMD ["./run"]