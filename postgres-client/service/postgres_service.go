package service

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mariajz/go-utils/postgres-client/model"
)

var DB *sqlx.DB

func InitPostgresDB(config model.Config) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName)

	var err error
	DB, err = sqlx.Open("postgres", dsn)

	if err != nil {
		fmt.Println("error in opening db connection:", err)
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println("error in ping:", err)
		panic(err)
	}
	log.Println("Database connection established")
}
