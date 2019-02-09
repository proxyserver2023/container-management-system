package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/alamin-mahamud/container-management-system/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

// Repository - ...
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

// VesselRepository - ...
type VesselRepository struct {
	vessels []*pb.Vessel
}

type service struct {
	repo Repository
}

// FindAvailable - ...
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if (spec.Capacity <= vessel.Capacity) && (spec.MaxWeight <= vessel.MaxWeight) {
			return vessel, nil
		}
	}
	return nil, errors.New("No Vessel found by that spec")
}

// FindAvailable - ...
func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {

	// Find the next available vessel
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	// Set the vessel as part of the response message type
	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{
			Id:        "ABC-001",
			Name:      "DEF-223-GHI-AJK",
			MaxWeight: 20000,
			Capacity:  500,
		},
	}

	repo := &VesselRepository{vessels}

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
