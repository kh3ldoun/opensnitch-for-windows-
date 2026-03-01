package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/evilsocket/opensnitch-windows/daemon/proto"
)

// GRPCServer manages the backend connection to the Python UI
type GRPCServer struct {
	pb.UnimplementedUIServer
}

func startGRPCServer(ctx context.Context) {
	// Typically on Windows, gRPC communicates over a local port or a named pipe.
	// For security and multi-user setups, named pipes are strongly preferred.
	// Example path: \\.\pipe\opensnitch
	pipePath := `\\.\pipe\opensnitch`

	// Create listener for named pipe
	listener, err := net.Listen("pipe", pipePath)
	if err != nil {
		log.Printf("Failed to listen on named pipe %s: %v", pipePath, err)
		// Fallback to TCP if pipe fails (e.g. perms)
		listener, err = net.Listen("tcp", "127.0.0.1:50051")
		if err != nil {
			log.Fatalf("Failed to listen on TCP fallback: %v", err)
		}
		log.Printf("Listening for gRPC on 127.0.0.1:50051 (TCP fallback)")
	} else {
		log.Printf("Listening for gRPC on named pipe %s", pipePath)
	}

	server := grpc.NewServer()
	pb.RegisterUIServer(server, &GRPCServer{})

	go func() {
		<-ctx.Done()
		log.Println("Stopping gRPC server...")
		server.GracefulStop()
	}()

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

// Ping implements pb.UIServer.Ping
func (s *GRPCServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{}, nil
}
