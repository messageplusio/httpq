package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/messgeplusio/httpq/pkg"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database successfully")

	err = pkg.Enqueue(db, "Hello, World!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Message enqueued successfully")

	message, err := pkg.Dequeue(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dequeued message:", message)
}
