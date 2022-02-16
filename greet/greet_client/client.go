package main

import (
	"context"
	"fmt"
	"grpc-go-sample/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main()  {
	fmt.Println("Welcome to Greet Client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	c := greetpb.NewGreetServiceClient(conn)
	//doUnary(c)
	//doServerStreaming(c)
	doClientStream(c)
}


func doUnary(c greetpb.GreetServiceClient)  {
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

func doServerStreaming(c greetpb.GreetServiceClient)  {
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
	for  {
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

func doClientStream(c greetpb.GreetServiceClient)  {
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName:  "Victor",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName:  "Cen man",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName:  "Dan fo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName:  "Qui",
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
