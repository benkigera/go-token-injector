package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Broker   string
	Port     int
	Topic    string
	Username string
	Password string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	Broker = os.Getenv("MQTT_BROKER")
	portStr := os.Getenv("MQTT_PORT")
	if portStr == "" {
		log.Fatal("MQTT_PORT not set in .env")
	}
	var err error
	Port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid MQTT_PORT: %v", err)
	}

	Username = os.Getenv("MQTT_USERNAME")
	Password = os.Getenv("MQTT_PASSWORD")
	Topic = os.Getenv("MQTT_TOPIC")
	if Topic == "" {
		log.Fatal("MQTT_TOPIC not set in .env")
	}
}