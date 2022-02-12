package main

import (
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
	fmt.Printf("Created Client: %f", c)
}
