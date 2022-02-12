package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-sample/greet/greetpb"
	"log"
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
	doUnary(c)
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
