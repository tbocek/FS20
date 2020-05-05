#Dockerfile
FROM golang:alpine AS builder
WORKDIR /build
COPY storage.pb.go server-stateless.go ./
RUN apk --no-cache add git wait4ports
RUN go get -d -v
RUN go build server-stateless.go storage.pb.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/server-stateless .
CMD ["wait4ports", "tcp:127.0.0.1:8082", "--", "./server-stateless", "-p", "8080", "-s", "8082"]
