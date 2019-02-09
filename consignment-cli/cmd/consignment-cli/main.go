package main

import (
	"context"
	"log"
	"os"

	"github.com/alamin-mahamud/container-management-system/consignment-cli/pkg/utl"
	pb "github.com/alamin-mahamud/container-management-system/consignment-service/proto/consignment"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/metadata"
)

const (
	defaultFileName = "consignment.json"
)

func main() {
	err := cmd.Init()

	if err != nil {
		log.Fatalf("Could not initialize the server -> %v", err)
	}

	// Create a new greeter client
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)

	file := defaultFileName
	var token string

	if len(os.Args) > 1 {
		file = os.Args[1]
		token = os.Args[2]
	}

	consignment, err := utl.ParseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	r, err := client.CreateConsignment(ctx, consignment)

	if err != nil {
		log.Fatalf("Could not Create consignment: %v", err)
	}

	log.Printf("Created %t", r.Created)

	getAll, err := client.GetConsignments(
		ctx,
		&pb.GetRequest{},
	)

	if err != nil {
		log.Fatalf("Could not list consignments - %v", err)
	}

	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
