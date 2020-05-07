//go build server-storage.go storage.pb.go
package main

import (
	context "context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	kvm         = make(map[string]string)
	optsStorage struct {
		Port int `short:"p" long:"port" description:"Port to listen on" required:"true"`
	}
)

type server struct{}

func (s server) GetKey(ctx context.Context, k *Key) (*Value, error) {
	v := kvm[k.Key]
	return &Value{Value: v}, nil
}

func (s server) PutKeyValue(ctx context.Context, kv *KeyValue) (*Empty, error) {
	kvm[kv.Key] = kv.Value
	return &Empty{}, nil
}

func main() {
	_, err := flags.NewParser(&optsStorage, flags.None).Parse()
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", optsStorage.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterStorageServer(grpcServer, &server{})
	log.Printf("Starting storage on port %v...", optsStorage.Port)
	log.Fatal(grpcServer.Serve(lis))
}
