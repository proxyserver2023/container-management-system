package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/alamin-mahamud/container-management-system/user-service/proto/user"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
)

func main() {
	cmd.Init()
	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)
	name := "Alamin Mahamud"
	email := "alamin.ineedahelp@gmail.com"
	password := "abc123"
	company := "ABC - Company"

	log.Println(name, email, password)
	fmt.Println("789")
	r, err := client.Create(context.TODO(), &pb.User{
		Name:     name,
		Email:    email,
		Password: password,
		Company:  company,
	})

	if err != nil {
		log.Fatalf("Couldn't create: %v", err)
	}

	log.Printf("Created: %s", r.User.Id)
	getAll, err := client.GetAll(context.Background(), &pb.Request{})

	if err != nil {
		log.Fatalf("Couldn't list users: %v", err)
	}

	for _, v := range getAll.Users {
		log.Println(v)
	}

	authResponse, err := client.Auth(context.TODO(), &pb.User{
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Could not authenticate user: %s error: %v\n", email, err)
	}

	log.Printf("Your access token is: %s\n", authResponse.Token)

	// let's just exit because
	os.Exit(0)
}
