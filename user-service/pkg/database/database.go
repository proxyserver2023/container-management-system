package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// CreateConnection - ...
func CreateConnection() (*gorm.DB, error) {
	host := "localhost"                           // os.Getenv("DB_HOST")
	databaseUser := "postgres"                    // os.Getenv("DB_USER")
	databaseName := "container_management_system" // os.Getenv("DB_NAME")
	databasePassword := "postgres-simple"         // os.Getenv("DB_PASSWORD")

	return gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s user=%s dbname=%s sslmode=disable password=%s",
			host, databaseUser, databaseName, databasePassword,
		),
	)
}
