package handler

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/alamin-mahamud/container-management-system/user-service/pkg/jwt"
	"github.com/alamin-mahamud/container-management-system/user-service/pkg/repository"
	pb "github.com/alamin-mahamud/container-management-system/user-service/proto/user"
)

// Service - ...
type Service struct {
	Repo         repository.Repository
	TokenService jwt.Authable
}

func (srv *Service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := srv.Repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (srv *Service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := srv.Repo.GetAll()
	if err != nil {
		return err
	}

	res.Users = users
	return nil
}

func (srv *Service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.Repo.GetByEmail(req.Email)

	log.Println(user)
	if err != nil {
		return err
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.TokenService.Encode(user)
	if err != nil {
		return err
	}

	res.Token = token
	return nil
}

func (srv *Service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	// generates a hashed version of our password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashedPass)
	if err := srv.Repo.Create(req); err != nil {
		return err
	}

	res.User = req
	return nil
}

func (srv *Service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	// Decode Token
	claims, err := srv.TokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	log.Println(claims)
	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	res.Valid = true
	return nil
}
