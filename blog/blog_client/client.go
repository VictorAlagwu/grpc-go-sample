package main

import (
	"context"
	"fmt"
	"grpc-go-sample/blog/blogpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Welcome to Blog Client")

	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Unable to connect to server: %v", err)
		}
	}(conn)

	c := blogpb.NewBlogServiceClient(conn)

	//handleCreateBlog(c)
	//handleReadBlog(c)
	//handleUpdateBlog(c)
	//handleDeleteBlog(c)
	handleFetchListData(c)
}

func handleCreateBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Mary",
			Title:    "Next Thing",
			Content:  "Yeah",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.CreateBlog(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline exceeded")
			} else {
				fmt.Printf("Unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Issue with CreateBlog rpc: %v", err)
		}
	}

	log.Printf("Response from CreateBlog: %v", res.Blog)
}

func handleReadBlog (c blogpb.BlogServiceClient) {
	req := &blogpb.ReadBlogRequest{
		BlogId: "62148e61c79d911ac8b9927f",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	res, err := c.ReadBlog(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline exceeded")
			} else {
				fmt.Printf("Unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Issue with ReadBlog rpc: %v", err)
		}
	}

	log.Printf("Response from ReadBlog: %v", res.Blog)
}

func handleUpdateBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id: "62148e61c79d911ac8b9927f",
			AuthorId: "Mary",
			Title:    "Next Thing",
			Content:  "Yeah",
		},
	}
	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Issue with UpdateBlog rpc: %v", err)
	}

	log.Printf("Response from UpdateBlog: %v", res.Blog)
}

func handleDeleteBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.DeleteBlogRequest{
		BlogId: "62148e61c79d911ac8b9927f",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	res, err := c.DeleteBlog(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline exceeded")
			} else {
				fmt.Printf("Unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Issue with DeleteBlog rpc: %v", err)
		}
	}

	log.Printf("Response from DeleteBlog: %v", res.Status)
}

func handleFetchListData(c blogpb.BlogServiceClient)  {
	resStream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Issue with ListBlog rpc: %v", err)
	}

	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}