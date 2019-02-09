package main

import (
	"fmt"
	"log"

	"github.com/alamin-mahamud/container-management-system/user-service/pkg/database"
	"github.com/alamin-mahamud/container-management-system/user-service/pkg/handler"
	"github.com/alamin-mahamud/container-management-system/user-service/pkg/jwt"
	"github.com/alamin-mahamud/container-management-system/user-service/pkg/repository"
	pb "github.com/alamin-mahamud/container-management-system/user-service/proto/user"
	micro "github.com/micro/go-micro"
)

func main() {
	// Creates a database connection and handles
	// closing it again before exit.
	db, err := database.CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	db.AutoMigrate(&pb.User{})
	repo := &repository.UserRepository{db}
	tokenService := &jwt.TokenService{repo}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Init will parse the cli flags
	srv.Init()

	// Register Handler
	pb.RegisterUserServiceHandler(srv.Server(), &handler.Service{repo, tokenService})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
