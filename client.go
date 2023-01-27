// function to call user grpc service

package user

import (
	"context"
	"log"

	pb "github.com/iamvasanth07/showcase/user/proto"
	"google.golang.org/grpc"
)

func callUserService() {
	// create a connection to the user service
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	// create a client
	c := pb.NewUserServiceClient(conn)
	// call the user service
	req := &pb.CreateUserRequest{
		User: &pb.User{
			Name:  "John Doe",
			Email: "",
			Phone: "1234567890",
		},
	}
	ctx := context.Background()
	res, err := c.Create(ctx, req)
	if err != nil {
		log.Fatalf("error while calling CreateUser RPC: %v", err)
	}
	log.Printf("Response from CreateUser: %v", res.User)
}
