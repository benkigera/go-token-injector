package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	injector_handlers "mqqt_go/api/injector/handlers"
	injector_routes "mqqt_go/api/injector/routes"
	injector_services "mqqt_go/api/injector/services"
	updater_handlers "mqqt_go/api/updater/handlers"
	updater_repositories "mqqt_go/api/updater/repositories"
	updater_routes "mqqt_go/api/updater/routes"
	updater_services "mqqt_go/api/updater/services"
	"mqqt_go/config"
	"mqqt_go/database"
	"mqqt_go/mqtt/client"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Printf("MQTT_PORT from environment: %s\n", os.Getenv("MQTT_PORT"))
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
	meterDataRepo := updater_repositories.NewMeterDataRepository("./latest_meter_data.txt")
	meterDataService := updater_services.NewMeterDataService(meterDataRepo)
	meterDataHandler := updater_handlers.NewMeterDataHandler(meterDataService)

	// Setup Injection API components
	injectionService := injector_services.NewInjectionService(client)
	injectionHandler := injector_handlers.NewInjectionHandler(injectionService)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	

	// Setup API routes
	updater_routes.SetupMeterDataRoutes(router, meterDataHandler)
	injector_routes.SetupInjectionRoutes(router, injectionHandler)

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