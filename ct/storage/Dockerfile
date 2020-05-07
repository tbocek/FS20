#Dockerfile
#to build image: docker build -t storage .
#to run: docker run -it --rm -p 8082:8080 storage
FROM golang:alpine AS builder
WORKDIR /build
COPY storage.pb.go server-storage.go ./
RUN apk --no-cache add git
RUN go get -d -v
RUN go build server-storage.go storage.pb.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/server-storage .
CMD ["./server-storage", "-p", "8080"]
