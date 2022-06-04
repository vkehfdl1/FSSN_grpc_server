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
	"strconv"
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

type clientStreamingServer struct {
	pb.UnimplementedClientStreamingServer
}

type serverStreamingServer struct {
	pb.UnimplementedServerStreamingServer
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

func (s *clientStreamingServer) GetServerResponse(stream pb.ClientStreaming_GetServerResponseServer) error {
	log.Printf("Server processing gRPC client-streaming.")
	var count int32 = 0
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Number{Value: count})
		}
		if err != nil {
			return err
		}
		if message != nil {
			count += 1
		}
	}
}

func (s *serverStreamingServer) GetServerResponse(num *pb.ServerNumber, stream pb.ServerStreaming_GetServerResponseServer) error {
	log.Printf("Server processing gRPC server-streaming %d", num.Value)
	for i := 0; i < 5; i++ {
		var message = &pb.ServerMessage{Message: "message #" + strconv.Itoa(i+1)}
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	serverName := flag.Int("server", 1, "which server to execute in this server")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Printf("server listening at %v", lis.Addr())
	if *serverName == 1 {
		// Example 1
		pb.RegisterMyServiceServer(s, &server{})
	} else if *serverName == 2 {
		pb.RegisterBidirectionalServer(s, &bidirectionalServer{})
	} else if *serverName == 3 {
		pb.RegisterClientStreamingServer(s, &clientStreamingServer{})
	} else if *serverName == 4 {
		pb.RegisterServerStreamingServer(s, &serverStreamingServer{})
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
