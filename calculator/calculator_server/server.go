package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-sample/calculator/calculatorpb"
	"log"
	"net"
	"strconv"
	"time"
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
