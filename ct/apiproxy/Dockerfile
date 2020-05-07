#Dockerfile
#to build image: docker build -t apiproxy .
#to run: docker run -it --rm -p 8081:8080 apiproxy
FROM golang:alpine AS builder
WORKDIR /build
COPY proxy.go .
RUN apk --no-cache add git
RUN go get -d -v
RUN go build proxy.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/proxy .
CMD ["./proxy"]
