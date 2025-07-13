package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin" // Added for Gin framework
	"mqqt_go/api/updater/handlers"
	"mqqt_go/api/updater/repositories"
	"mqqt_go/api/updater/routes"
	"mqqt_go/api/updater/services"
	"mqqt_go/config"
	"mqqt_go/database"
	"mqqt_go/mqtt/client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Initialize database
	database.InitDB(
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	defer database.DB.Close()
	// Enable MQTT client logging for debugging
	mqtt.DEBUG = log.New(os.Stdout, "[MQTT-DEBUG] ", log.LstdFlags)
	mqtt.ERROR = log.New(os.Stdout, "[MQTT-ERROR] ", log.LstdFlags)
	
	fmt.Println("Starting MQTT client...")
	
	client := client.CreateMQTTClient()
	
	// Connect with retry logic
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("Failed to connect: %v, retrying in 5 seconds...", token.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	// Setup API components
	meterDataRepo := repositories.NewMeterDataRepository("./latest_meter_data.txt")
	meterDataService := services.NewMeterDataService(meterDataRepo)
	meterDataHandler := handlers.NewMeterDataHandler(meterDataService)

	// Initialize Gin router
	router := gin.Default()

	// Setup API routes
	routes.SetupMeterDataRoutes(router, meterDataHandler)

	// Start HTTP server in a new goroutine
	go func() {
		log.Printf("HTTP server starting on port %s", config.APIPort)
		if err := router.Run(":" + config.APIPort); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()
	
	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	// Keep the main goroutine running
	go func() {
		<-c
		fmt.Println("\nShutting down...")
		client.Disconnect(250)
		os.Exit(0)
	}()
	
	// Monitor connection status
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if client.IsConnected() {
				fmt.Printf("[%s] Connection status: OK\n", time.Now().Format("2006-01-02 15:04:05"))
			} else {
				fmt.Printf("[%s] Connection status: DISCONNECTED\n", time.Now().Format("2006-01-02 15:04:05"))
			}
		}
	}
}