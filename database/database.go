package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(host, port, user, password, dbname string) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Ping the database to ensure connection is established
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS meter_data (
		name TEXT PRIMARY KEY,
		value TEXT,
		unit TEXT,
		last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		reading_timestamp TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS mqtt_messages (
		id SERIAL PRIMARY KEY,
		topic TEXT NOT NULL,
		payload TEXT NOT NULL,
		received_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Table 'meter_data' ensured.")
}

func InsertOrUpdateMeterData(name string, value interface{}, unit string, readingTimestamp time.Time) error {
	stmt, err := DB.Prepare(`
	INSERT INTO meter_data(name, value, unit, reading_timestamp)
	VALUES($1, $2, $3, $4)
	ON CONFLICT(name) DO UPDATE SET
		value = EXCLUDED.value,
		unit = EXCLUDED.unit,
		last_updated = CURRENT_TIMESTAMP,
		reading_timestamp = EXCLUDED.reading_timestamp;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, fmt.Sprintf("%v", value), unit, readingTimestamp)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func GetMeterData(name string) (string, string, error) {
	var value, unit string
	row := DB.QueryRow("SELECT value, unit FROM meter_data WHERE name = $1", name)
	if err := row.Scan(&value, &unit); err != nil {
		if err == sql.ErrNoRows {
			return "", "", nil // No rows found, return empty strings and no error
		}
		return "", "", fmt.Errorf("failed to query meter data: %w", err)
	}
	return value, unit, nil
}