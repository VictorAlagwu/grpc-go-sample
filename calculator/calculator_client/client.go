package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-sample/calculator/calculatorpb"
	"io"
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

	//doUnary(c)
	doServerStream(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient)  {
	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 10,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Issue with calculator rpc: %v", err)
	}

	log.Printf("Response from Calculator: %v", res.Result)
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Value: 210,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Issue with Prime number decomposition rpc: %v", err)
	}

	for  {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %v", msg.GetResult())
	}
}