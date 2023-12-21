package main

import (
	"context"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50051"
	userID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	log.Printf(color.RedString("User info:\n", color.GreenString("%+v", r.GetUser())))
}
