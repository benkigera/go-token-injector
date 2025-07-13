package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"mqqt_go/api/updater/models"
	"mqqt_go/api/updater/repositories"
	"mqqt_go/config"
	"mqqt_go/database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[%s] DEBUG: messagePubHandler called for topic: %s\n", 
		time.Now().Format("2006-01-02 15:04:05"), msg.Topic())
	fmt.Printf("[%s] Received message from topic: %s\n", 
		time.Now().Format("2006-01-02 15:04:05"), msg.Topic())

	// Initialize MQTT message repository
	mqttMsgRepo := repositories.NewMQTTMessageRepository(database.DB)

	// Create MQTTMessage model
	mqttMessage := &models.MQTTMessage{
		Topic:     msg.Topic(),
		Payload:   string(msg.Payload()),
		ReceivedAt: time.Now(),
	}

	// Save raw MQTT message to database
	if err := mqttMsgRepo.SaveMQTTMessage(context.Background(), mqttMessage); err != nil {
		log.Printf("[%s] ERROR: Failed to save raw MQTT message to DB: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		fmt.Printf("[%s] DEBUG: Raw MQTT message saved to DB with ID: %d\n", time.Now().Format("2006-01-02 15:04:05"), mqttMessage.ID)
	}

	// Unmarshal the JSON payload
	var data []map[string]interface{} // Declare data here
	fmt.Printf("[%s] Received JSON payload: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg.Payload())
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}
	fmt.Printf("[%s] DEBUG: JSON unmarshaled successfully. Items found: %d\n", time.Now().Format("2006-01-02 15:04:05"), len(data))

	// Generate a single timestamp for this reading
	readingTimestamp := time.Now()

	// Store data in the database
	for _, item := range data {
		name, _ := item["n"].(string)
		value := item["v"]
		unit, _ := item["u"].(string)

		fmt.Printf("[%s] DEBUG: Attempting to insert/update: Name=%s, Value=%v, Unit=%s\n", time.Now().Format("2006-01-02 15:04:05"), name, value, unit)
		if err := database.InsertOrUpdateMeterData(name, value, unit, readingTimestamp); err != nil {
			log.Printf("[%s] ERROR: Failed to insert/update data for %s: %v\n", time.Now().Format("2006-01-02 15:04:05"), name, err)
		} else {
			fmt.Printf("[%s] DEBUG: Successfully inserted/updated data for %s\n", time.Now().Format("2006-01-02 15:04:05"), name)
		}
	}

	// Read AvailableCredit from the database
	availableCreditValue, availableCreditUnit, err := database.GetMeterData("AvailableCredit")
	if err != nil {
		log.Printf("Error getting AvailableCredit from DB: %v", err)
		return
	}

	// Prepare content for the output file
	var outputContent string
	outputContent += fmt.Sprintf("--- Update at %s ---\n", time.Now().Format("2006-01-02 15:04:05"))

	if availableCreditValue != "" {
		outputContent += fmt.Sprintf("Available Credit: %s %s\n", availableCreditValue, availableCreditUnit)
	} else {
		outputContent += "Available Credit: Not available\n"
	}

	// Define the output filename
	filename := "latest_meter_data.txt"

	// Write the formatted content to the file, overwriting if it exists
	if err := os.WriteFile(filename, []byte(outputContent), 0644); err != nil {
		log.Printf("Error writing message to file %s: %v", filename, err)
		return
	}

	fmt.Printf("[%s] Message saved to %s\n", time.Now().Format("2006-01-02 15:04:05"), filename)
}

var ConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Printf("[%s] Connected to MQTT broker\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// Subscribe with retry logic
	if token := client.Subscribe(config.Topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe: %v", token.Error())
		return
	}
	fmt.Printf("[%s] Subscribed to topic: %s\n", time.Now().Format("2006-01-02 15:04:05"), config.Topic)
}

var ConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("[%s] Connection lost: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
}

var ReconnectHandler mqtt.ReconnectHandler = func(client mqtt.Client, opts *mqtt.ClientOptions) {
	fmt.Printf("[%s] Attempting to reconnect...\n", time.Now().Format("2006-01-02 15:04:05"))
}