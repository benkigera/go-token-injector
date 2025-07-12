package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[%s] DEBUG: messagePubHandler called for topic: %s\n", 
		time.Now().Format("2006-01-02 15:04:05"), msg.Topic())
	fmt.Printf("[%s] Received message from topic: %s\n", 
		time.Now().Format("2006-01-02 15:04:05"), msg.Topic())

	// Unmarshal the JSON payload
	var data []map[string]interface{}
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	// Prepare content for the output file
	var outputContent string
	outputContent += fmt.Sprintf("--- Update at %s ---\n", time.Now().Format("2006-01-02 15:04:05"))

	for _, item := range data {
		name, _ := item["n"].(string)
		value := item["v"]
		unit, hasUnit := item["u"].(string)

		if hasUnit {
			outputContent += fmt.Sprintf("Name: %s, Value: %v, Unit: %s\n", name, value, unit)
		} else {
			outputContent += fmt.Sprintf("Name: %s, Value: %v\n", name, value)
		}
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

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Printf("[%s] Connected to MQTT broker\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// Subscribe with retry logic
	if token := client.Subscribe(Topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe: %v", token.Error())
		return
	}
	fmt.Printf("[%s] Subscribed to topic: %s\n", time.Now().Format("2006-01-02 15:04:05"), Topic)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("[%s] Connection lost: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
}

var reconnectHandler mqtt.ReconnectHandler = func(client mqtt.Client, opts *mqtt.ClientOptions) {
	fmt.Printf("[%s] Attempting to reconnect...\n", time.Now().Format("2006-01-02 15:04:05"))
}