#Dockerfile
#to build image: docker build -t service .
#to run: docker run -it --rm -p 8085:8080 service
FROM golang:alpine AS builder
WORKDIR /build
COPY storage.pb.go server-stateless.go ./
RUN apk --no-cache add git
RUN go get -d -v
RUN go build server-stateless.go storage.pb.go

FROM alpine:latest
WORKDIR /root/
RUN apk --no-cache add wait4ports
COPY --from=builder /build/server-stateless .
CMD ["./server-stateless", "-p", "8080", "-s", "8082"]