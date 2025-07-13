package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"mqqt_go/api/updater/models"
)

type MQTTMessageRepository interface {
	SaveMQTTMessage(ctx context.Context, msg *models.MQTTMessage) error
}

type mqttMessageRepository struct {
	db *sql.DB
}

func NewMQTTMessageRepository(db *sql.DB) MQTTMessageRepository {
	return &mqttMessageRepository{db: db}
}

func (r *mqttMessageRepository) SaveMQTTMessage(ctx context.Context, msg *models.MQTTMessage) error {
	query := `
		INSERT INTO mqtt_messages (topic, payload, received_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query, msg.Topic, msg.Payload, msg.ReceivedAt).Scan(&msg.ID)
	if err != nil {
		return fmt.Errorf("failed to insert MQTT message: %w", err)
	}
	return nil
}
