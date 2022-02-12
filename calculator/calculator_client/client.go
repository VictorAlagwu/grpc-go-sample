package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-sample/calculator/calculatorpb"
	"log"
)

func main()  {
	fmt.Println("Welcome to gRPC calculator")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	c := calculatorpb.NewCalculatorServiceClient(conn)

	doUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient)  {
	req := &calculatorpb.CalculatorRequest{
		Calculating: &calculatorpb.Calculating{
		FirstNumber:  3,
		SecondNumber: 10,
	}}

	res, err := c.Calculator(context.Background(), req)
	if err != nil {
		log.Fatalf("Issue with calculator rpc: %v", err)
	}

	log.Printf("Response from Calculator: %v", res.Result)
}
