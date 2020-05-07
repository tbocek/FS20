#Dockerfile
FROM golang:alpine AS builder
WORKDIR /build
COPY auth.go ./
RUN apk --no-cache add git
RUN go get -d -v
RUN go build auth.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/auth .
CMD ["./auth", "-p", "8080"]
