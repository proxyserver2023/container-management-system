package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	pb "github.com/alamin-mahamud/container-management-system/consignment-service/proto/consignment"
	vesselProto "github.com/alamin-mahamud/container-management-system/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
)

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo         Repository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	vesselResopnse, err := s.vesselClient.FindAvailable(
		context.Background(),
		&vesselProto.Specification{
			MaxWeight: req.Weight,
			Capacity:  int32(len(req.Containers)),
		},
	)

	log.Printf("Found vessel: %s \n", vesselResopnse.Vessel.Name)

	if err != nil {
		return err
	}

	req.VesselId = vesselResopnse.Vessel.Id

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}

// GetConsignment - returns all the created consignments
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

const (
	defaultHost = "localhost:27017"
)

func main() {

	// Database host from the environment variable
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}

	sess, err := CreateSession(host)

	// Mgo creates a 'master' session, we need to end that session
	// before the main function closes.
	defer session.Close()

	if err != nil {
		// We're wrapping the error returned from our CreateSession
		// here to add some context to the error.
		log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		micro.WrapHandler(AuthWrapper),
	)

	repo := &ConsignmentRepository{}

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// Init will parse the command line flags
	srv.Init()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

}

// AuthWrapper - ...
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}
		token := meta["Token"]
		log.Println("Authenticating with token: ", token)

		// Auth here
		authClient := userService.NewUserServiceClient("go.micro.srv.user", client.DefaultClient)
		authResp, err := authClient.ValidateToken(ctx, &userService.Token{
			Token: token,
		})

		log.Println("Auth resp:", authResp)
		log.Println("Err:", err)
		if err != nil {
			return err
		}

		err = fn(ctx, req, resp)
		return err
	}
}
