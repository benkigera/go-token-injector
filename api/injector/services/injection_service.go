package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"mqqt_go/api/injector/models"
	"mqqt_go/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type InjectionService interface {
	InjectToken(ctx context.Context, req *models.InjectionRequest) (models.MQTTInjectionResponsePayload, error)
}

type injectionService struct {
	mqttClient mqtt.Client
}

func NewInjectionService(mqttClient mqtt.Client) InjectionService {
	return &injectionService{mqttClient: mqttClient}
}

func (s *injectionService) InjectToken(ctx context.Context, req *models.InjectionRequest) (models.MQTTInjectionResponsePayload, error) {
	responseTopic := fmt.Sprintf("%s/%s", config.MQTTInjectionResponseTopicBase, req.MeterNumber)
	
	// Channel to receive MQTT response
	responseChan := make(chan models.MQTTInjectionResponsePayload, 1)

	// Dynamically subscribe to the response topic
	token := s.mqttClient.Subscribe(responseTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received response on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
		var mqttResponse models.MQTTInjectionResponsePayload
		if err := json.Unmarshal(msg.Payload(), &mqttResponse); err != nil {
			log.Printf("Error unmarshaling MQTT response: %v", err)
			return
		}
		select {
		case responseChan <- mqttResponse:
			// Message sent successfully
		case <-time.After(1 * time.Second): // Small timeout to prevent blocking
			log.Println("Timeout sending MQTT response to channel")
		}
	})
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to subscribe to response topic %s: %w", responseTopic, token.Error())
	}
	log.Printf("Subscribed to dynamic response topic: %s\n", responseTopic)

	defer func() {
		// Unsubscribe after processing or timeout
		unsubToken := s.mqttClient.Unsubscribe(responseTopic)
		if unsubToken.Wait() && unsubToken.Error() != nil {
			log.Printf("Error unsubscribing from topic %s: %v", responseTopic, unsubToken.Error())
		}
		log.Printf("Unsubscribed from dynamic response topic: %s\n", responseTopic)
	}()

	// Prepare payload for external API
	externalAPIPayload := models.ExternalAPIPayload{
		{
			N: "1P-Energy-SerialNumber",
			V: req.MeterNumber,
			T: time.Now().Unix(),
		},
		{
			N: "credit-token-injection",
			V: req.CreditToken,
			T: time.Now().Unix(),
		},
	}
	payloadBytes, err := json.Marshal(externalAPIPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal external API payload: %w", err)
	}

	// Make HTTP POST request to external API
	apiURL := fmt.Sprintf("%s/%s", config.ExternalAPIBaseURL, req.MeterNumber)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", config.ExternalAPIAuthToken)

	client := &http.Client{Timeout: 10 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to external API: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("external API returned non-success status: %d, body: %s", httpResp.StatusCode, string(bodyBytes))
	}

	log.Printf("Successfully sent injection request for meter %s\n", req.MeterNumber)

	// Wait for MQTT response with a timeout
	select {
	case mqttResponse := <-responseChan:
		return mqttResponse, nil
	case <-time.After(30 * time.Second): // Configurable timeout
		return nil, fmt.Errorf("timeout waiting for MQTT response for meter %s", req.MeterNumber)
	}
}
