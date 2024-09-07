package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	// TODO 환경 변수 분리
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE")
	port := os.Getenv("DATABASE_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		os.Exit(1)
	}

	fmt.Println("Database connected successfully")

	return db
}
