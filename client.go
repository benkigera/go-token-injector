package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func createMQTTClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	
	// Connection settings
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", Broker, Port))
	opts.SetClientID(fmt.Sprintf("go_mqtt_client_%d", time.Now().Unix()))
	opts.SetUsername(Username)
	opts.SetPassword(Password)
	
	// Keep-alive and timeout settings
	opts.SetKeepAlive(30 * time.Second)          // Reduced from 60s for faster detection
	opts.SetPingTimeout(10 * time.Second)        // Ping timeout
	opts.SetConnectTimeout(30 * time.Second)     // Connection timeout
	opts.SetWriteTimeout(10 * time.Second)       // Write timeout
	
	// Auto-reconnect settings
	opts.SetAutoReconnect(true)                  // Enable auto-reconnect
	opts.SetMaxReconnectInterval(30 * time.Second) // Max time between reconnect attempts
	opts.SetConnectRetryInterval(5 * time.Second)  // Initial retry interval
	opts.SetConnectRetry(true)                   // Enable connect retry
	
	// Clean session for reliable message delivery
	opts.SetCleanSession(true)
	
	// Quality of Service settings
	opts.SetOrderMatters(false)                  // Allow out-of-order message processing
	opts.SetResumeSubs(true)                     // Resume subscriptions on reconnect
	
	// Message handling
	opts.SetDefaultPublishHandler(messagePubHandler)
	
	// Event handlers
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.OnReconnecting = reconnectHandler
	
	// Create client
	client := mqtt.NewClient(opts)
	return client
}