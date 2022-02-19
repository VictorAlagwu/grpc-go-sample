package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-sample/calculator/calculatorpb"
	"io"
	"log"
	"time"
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
	//doServerStream(c)
	//doClientStream(c)
	doBiDiStreaming(c)
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

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	requests := []*calculatorpb.ComputeAverageRequest{
		{
			Value: 1,
		}, {
			Value: 2,
		}, {
			Value: 3,
		}, {
			Value: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Issue with compute average")
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error while reading single stream: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving stream: %v", err)
	}
	fmt.Printf("Compute Average Response: %v\n", res)
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Send Client Request on Bi-directional Streaming")
	requests := []*calculatorpb.FindMaximumRequest{
		{
			Value: 1,
		}, {
			Value: 5,
		}, {
			Value: 3,
		}, {
			Value: 6,
		}, {
			Value: 2,
		}, {
			Value: 20,
		},
	}
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error processing bi-di request: %v", err)
		return
	}

	processor := make(chan struct{})

	// Send Request to client
	go func() {
		for _, req := range requests {
			fmt.Printf("Processing message: %v\n", req)
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

	// Receive response
	go func() {
		for  {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error when receiving bi-dr response: %v", err)
			}
			fmt.Printf("Value: %v\n", res.GetValue())
		}
		close(processor)
	}()
	// Blocks until channel is resolved
	<-processor
}