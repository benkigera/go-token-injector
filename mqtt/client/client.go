package client

import (
	"fmt"
	"time"

	"mqqt_go/config"
	"mqqt_go/mqtt/handlers"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func CreateMQTTClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	
	// Connection settings
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Broker, config.Port))
	opts.SetClientID(fmt.Sprintf("go_mqtt_client_%d", time.Now().Unix()))
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	
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
	opts.SetDefaultPublishHandler(handlers.MessagePubHandler)
	
	// Event handlers
	opts.OnConnect = handlers.ConnectHandler
	opts.OnConnectionLost = handlers.ConnectLostHandler
	opts.OnReconnecting = handlers.ReconnectHandler
	
	// Create client
	mqttClient = mqtt.NewClient(opts)
	return mqttClient
}

func GetMQTTClient() mqtt.Client {
	return mqttClient
}