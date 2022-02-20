package main

import (
	"context"
	"fmt"
	"grpc-go-sample/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Invoking Greet Function %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("Invoking Greet with Deadline Function %v", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client canceled the request")
			return nil, status.Errorf(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 5; i++ {
		result := "Name: " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("Invoking LongGreet function with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//End of stream request
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error reading client stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Sending request " + firstName + "| "
	}

}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("Starting GreetEveryone service\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("End of request")
		}
		if err != nil {
			log.Fatalf("Error reading bi-directional stream: #{err}")
		}
		firstName := req.GetGreeting().GetFirstName()

		sendStreamErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Hello " + firstName + " \n",
		})
		if sendStreamErr != nil {
			return sendStreamErr
		}
	}
}

func main() {
	fmt.Println("In the beginning of GRPC")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Listener error: %v", err)
	}
	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslError := credentials.NewServerTLSFromFile(certFile, keyFile)

		if sslError != nil {
			log.Fatalf("Failed to load certificates: %v", sslError)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
