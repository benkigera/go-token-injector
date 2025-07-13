package models

import (
	"time"
)

// MQTTMessage represents a raw MQTT message received by the client.
type MQTTMessage struct {
	ID        int64     `json:"id" db:"id"`
	Topic     string    `json:"topic" db:"topic"`
	Payload   string    `json:"payload" db:"payload"` // Store raw JSON as string
	ReceivedAt time.Time `json:"received_at" db:"received_at"`
}
