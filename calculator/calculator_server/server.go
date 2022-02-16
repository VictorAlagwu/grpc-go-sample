package main

import (
	"context"
	"fmt"
	"grpc-go-sample/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("A request has been made %v", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	result := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(
	req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	value := req.GetValue()
	var evenNumber int32
	evenNumber = 2

	for value > 1 {
		if value % evenNumber == 0 {
			result := strconv.Itoa(int(evenNumber))
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: result,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			time.Sleep(1000 * time.Millisecond)
			value = value / evenNumber
		} else {
			evenNumber = evenNumber + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Invoking Compute Average\n")
	var result float64
	totalSum := 0
	result = 0
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result = float64(totalSum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error reading client stream: %v", err)
		}
		count++
		value := req.GetValue()
		totalSum += int(value)
	}
	return nil
}

func main() {
	fmt.Println("Thus said the programmer, go forth and render gRPC")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Listener error: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
