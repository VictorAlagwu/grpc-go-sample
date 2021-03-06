package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-go-sample/calculator/calculatorpb"
	"io"
	"log"
	"math"
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
		if value%evenNumber == 0 {
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

}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("Starting Find Maximum Request")
	var maximumValue int32
	maximumValue = 0
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			log.Fatalf("End of request")
		}

		if err != nil {
			log.Fatalf("Error reading bi-directional stream: #{err}")
		}

		value := req.Value
		if value > maximumValue {
			maximumValue = value
			sendStreamErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Value: maximumValue,
			})
			if sendStreamErr != nil {
				return sendStreamErr
			}
		}

	}
}

func (*server)SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Square root RPC request received")
	number := req.GetValue()
	
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number"),
			)
	}
	res := &calculatorpb.SquareRootResponse{
		Root: math.Sqrt(float64(number)),
	}
	return res, nil
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
