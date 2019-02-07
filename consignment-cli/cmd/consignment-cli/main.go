package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/alamin-mahamud/container-management-system/consignment-service/proto/consignment"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
)

const (
	defaultFileName = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, err
	}
	return consignment, err
}

func main() {
	err := cmd.Init()

	if err != nil {
		log.Fatalf("Could not initialize the server -> %v", err)
	}

	// Create a new greeter client
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)
	file := defaultFileName
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(
		context.Background(),
		consignment,
	)

	if err != nil {
		log.Fatalf("Could not Create consignment: %v", err)
	}

	log.Printf("Created %t", r.Created)

	getAll, err := client.GetConsignments(
		context.Background(),
		&pb.GetRequest{},
	)

	if err != nil {
		log.Fatalf("Could not list consignments - %v", err)
	}

	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
