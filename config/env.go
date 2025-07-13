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

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	APIPort string
	ExternalAPIBaseURL string
	ExternalAPIAuthToken string
	MQTTInjectionResponseTopicBase string
)

func init() {
	// Load .env file only if not in production environment
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v (This is expected in production environments)", err)
		}
	}

	Broker = os.Getenv("MQTT_BROKER")
	portStr := os.Getenv("MQTT_PORT")
	if portStr == "" {
		log.Fatal("MQTT_PORT not set")
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
		log.Fatal("MQTT_TOPIC not set")
	}

	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")

	APIPort = os.Getenv("API_PORT")
	if APIPort == "" {
		log.Fatal("API_PORT not set")
	}

	ExternalAPIBaseURL = os.Getenv("EXTERNAL_API_BASE_URL")
	if ExternalAPIBaseURL == "" {
		log.Fatal("EXTERNAL_API_BASE_URL not set")
	}

	ExternalAPIAuthToken = os.Getenv("EXTERNAL_API_AUTH_TOKEN")
	if ExternalAPIAuthToken == "" {
		log.Fatal("EXTERNAL_API_AUTH_TOKEN not set")
	}

	MQTTInjectionResponseTopicBase = os.Getenv("MQTT_INJECTION_RESPONSE_TOPIC_BASE")
	if MQTTInjectionResponseTopicBase == "" {
		log.Fatal("MQTT_INJECTION_RESPONSE_TOPIC_BASE not set")
	}
}