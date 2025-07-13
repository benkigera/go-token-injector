package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load("/Users/mac/Documents/work/mqqt_go/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get DB connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure connection is established
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection established.")

	rows, err := db.Query("SELECT id, topic, payload, received_at FROM mqtt_messages ORDER BY received_at DESC")
	if err != nil {
		log.Fatalf("Failed to query mqtt_messages: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n--- MQTT Messages ---")
	for rows.Next() {
		var id int64
		var topic string
		var payload string
		var receivedAt time.Time

		if err := rows.Scan(&id, &topic, &payload, &receivedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		fmt.Printf("ID: %d\n", id)
		fmt.Printf("Topic: %s\n", topic)
		fmt.Printf("Received At: %s\n", receivedAt.Format(time.RFC3339))
		fmt.Printf("Payload:\n%s\n", payload)
		fmt.Println("---------------------")
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
}