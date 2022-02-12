package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-sample/calculator/calculatorpb"
	"log"
	"net"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Calculator(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("A request has been made %v", req)
	firstNumber := req.GetCalculating().GetFirstNumber()
	secondNumber := req.GetCalculating().GetSecondNumber()
	result := firstNumber + secondNumber
	res := &calculatorpb.CalculatorResponse{
		Result: result,
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
