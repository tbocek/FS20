#Dockerfile
FROM golang:alpine AS builder
WORKDIR /build
COPY server-stateful.go .
RUN apk --no-cache add git
RUN go get -d -v
RUN go build server-stateful.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/server-stateful .
CMD ["./server-stateful", "-p", "8081"]
