package service

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mariajz/go-utils/postgres-client/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgresDB(config model.Config) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error in opening db connection:", err)
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("error in conversion:", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		fmt.Println("error in ping:", err)
		panic(err)
	}
	log.Println("Database connection established")
}
