package main

import (
	"context"
	"flag"
	"fmt"
	pb "fssn_grpc/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var (
	port = flag.Int("port", 54321, "The server port")
)

type server struct {
	pb.UnimplementedMyServiceServer
}

type bidirectionalServer struct {
	pb.UnimplementedBidirectionalServer
}

func (s *server) MyFunction(ctx context.Context, in *pb.MyNumber) (*pb.MyNumber, error) {
	log.Printf("Received: %d", in.GetValue())
	return &pb.MyNumber{Value: in.GetValue() * in.GetValue()}, nil
}

func (s *bidirectionalServer) GetServerResponse(stream pb.Bidirectional_GetServerResponseServer) error {
	log.Printf("Server processing gRPC bidirectional streaming.")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		resp := pb.Message{Message: in.Message}
		if err := stream.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		log.Printf("send message=%s", resp.Message)
	}
}

func main() {
	serverName := flag.Int("server", 1, "which server to execute in this server")
	flag.Parse()

	if *serverName == 1 {
		// Example 1
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterMyServiceServer(s, &server{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	} else if *serverName == 2 {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterBidirectionalServer(s, &bidirectionalServer{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
