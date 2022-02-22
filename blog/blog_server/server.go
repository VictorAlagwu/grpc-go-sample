package main

import (
	"context"
	"fmt"
	"grpc-go-sample/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)
var collection *mongo.Collection

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID   primitive.ObjectID   `bson:"_id,omitempty"`
	AuthorID string  `bson:"author_id"`
	Content string `bson:"content"`
	Title string `bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Create Blog request")
	blog := req.GetBlog()

	payload := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}
	one, err := collection.InsertOne(context.Background(), payload)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
			)
	}

	oid, ok := one.InsertedID.(primitive.ObjectID)
	if ! ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID: %v", err),
			)
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
			Content: blog.GetContent(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
		},
	},nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error){
	fmt.Println("Read Blog request")
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID: %v", err),
		)
	}

	data := &blogItem{}

	filter := primitive.M{"_id": oid}

	one := collection.FindOne(context.Background(), filter)
	
	if err := one.Decode(data); err != nil{
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Not Found error: %v", err),
		)
	}


	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id: data.ID.Hex(),
			Content: data.Content,
			AuthorId: data.AuthorID,
			Title: data.Title,
		},
	},nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Update Blog request")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID: %v", err),
		)
	}

	data := &blogItem{}

	filter := primitive.M{"_id": oid}
	one := collection.FindOne(context.Background(), filter)

	if err := one.Decode(data); err != nil{
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Not Found error: %v", err),
		)
	}
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unable to update object: %v", err),
		)
	}
	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id: data.ID.Hex(),
			Content: data.Content,
			AuthorId: data.AuthorID,
			Title: data.Title,
		},
	}, nil
}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Delete Blog request")
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID: %v", err),
		)
	}
	filter := primitive.M{"_id": oid}

	res, deleteErr := collection.DeleteOne(context.Background(), filter)

	if deleteErr != nil{
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Unable to delete: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog in MongoDB: %v", err),
			)
	}
	return &blogpb.DeleteBlogResponse{Status: "success"}, nil
}

func (*server) ListBlog(_ *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer)  error {
	res, err := collection.Find(context.Background(), primitive.D{})

	if err != nil{
		return status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Unable to Fetch data: %v", err),
		)
	}
	defer func(res *mongo.Cursor, ctx context.Context) {
		err := res.Close(ctx)
		if err != nil {

		}
	}(res, context.Background())
	for res.Next(context.Background()) {
		data := &blogItem{}
		err := res.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Unable to decode data: %v", err),
			)
		}
		streamErr := stream.Send(&blogpb.ListBlogResponse{
			Blog: &blogpb.Blog{
				Id: data.ID.Hex(),
				Content: data.Content,
				AuthorId: data.AuthorID,
				Title: data.Title,
			},
		})
		if streamErr != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Streaming error: %v", streamErr),
			)
		}
	}

	if err := res.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Starting Blog GRPC server")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("localdb").Collection("blog")



	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Listener error: %v", err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Will wait for control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	listenerErr := lis.Close()
	if listenerErr != nil {
		log.Fatalf("Listener error: %v", err)
		return
	}
	fmt.Println("Ending...")
}
