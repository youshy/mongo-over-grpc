package main

import (
	"context"
	"io"
	"log"

	"github.com/mongo-over-grpc/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create Blog
	log.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Artur",
		Title:    "Testing Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	log.Printf("Blog has been created:\n%v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// read Blog
	log.Println("Reading the blog")

	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if err != nil {
		log.Printf("Error happened while reading: %v\n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, err := c.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		log.Printf("Error happened while reading: %v\n", err)
	}

	log.Printf("Blog was read:\n%v\n", readBlogRes)

	// update Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		log.Printf("Error happened while updating: %v\n", err)
	}
	log.Printf("Blog was updated:\n%v\n", updateRes)

	// delete Blog
	deleteRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if err != nil {
		log.Printf("Error happened while deleting: %v\n", err)
	}
	log.Printf("Blog was deleted:\n%v\n", deleteRes)

	// list Blogs

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		log.Println(res.GetBlog())
	}
}
