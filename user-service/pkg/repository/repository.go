package repository

import (
	pb "github.com/alamin-mahamud/container-management-system/user-service/proto/user"
	"github.com/jinzhu/gorm"
)

// Repository - ...
type Repository interface {
	GetAll() ([]*pb.User, error)
	Get(id string) (*pb.User, error)
	Create(user *pb.User) error
	GetByEmail(email string) (*pb.User, error)
}

// UserRepository - ...
type UserRepository struct {
	Db *gorm.DB
}

// GetAll - ...
func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	if err := repo.Db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Get - ...
func (repo *UserRepository) Get(id string) (*pb.User, error) {
	var user *pb.User
	if err := repo.Db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail - ...
func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	var user *pb.User
	if err := repo.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Create - ...
func (repo *UserRepository) Create(user *pb.User) error {
	if err := repo.Db.Create(user).Error; err != nil {
		return err
	}

	return nil
}
