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

func (*server) GreetManyTimes( req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
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
	for  {
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

func main()  {
	fmt.Println("In the beginning of GRPC")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Listener error: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
