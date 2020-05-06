package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/mongo-over-grpc/blog/blogpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	collection *mongo.Collection
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Blog Service Started\n")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := &server{}
	server.collection = MongoSetup()

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, server)

	go func() {
		log.Printf("Starting Blog server\n")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// wait for signal to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	log.Printf("Stopping server\n")
	s.Stop()
	lis.Close()
}

func MongoSetup() *mongo.Collection {
	// TODO: env key
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("MongoDB ready\n")

	return client.Database("mydb").Collection("blog")
}
