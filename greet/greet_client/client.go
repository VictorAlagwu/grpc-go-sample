package main

import (
	"context"
	"fmt"
	"grpc-go-sample/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Welcome to Greet Client")

	certFile := "ssl/ca.crt" // Certificate Authority Trust certificate
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	if sslErr != nil {
		log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
		return
	}
	opts := grpc.WithTransportCredentials(creds)

	conn, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Unable to connect to server: %v", err)
		}
	}(conn)

	c := greetpb.NewGreetServiceClient(conn)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStream(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Victor",
			SecondName: "A",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Issue with Greet rpc: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Victor",
			SecondName: "A",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Issue with Greet rpc: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStream(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Victor",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Cen man",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dan fo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Qui",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Issue with LongGreet rpc: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving streamkm: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC....")
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Victor",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Cen man",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dan fo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Qui",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	// Make channel
	waitc := make(chan struct{})
	// Send messages to client (Go routine)
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			err := stream.Send(req)
			if err != nil {
				return
			}
			time.Sleep(1000 * time.Millisecond)
		}
		err := stream.CloseSend()
		if err != nil {
			return
		}
	}()

	// Receive messages from client (Go routine)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
			}
			fmt.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()

	// Block until everything is done
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Victor",
			SecondName: "A",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline exceeded")
			} else {
				fmt.Printf("Unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Issue with GreetWithDeadline rpc: %v", err)
		}
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
