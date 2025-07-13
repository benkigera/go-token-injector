package models

type InjectionRequest struct {
	MeterNumber string `json:"meter_number" binding:"required"`
	CreditToken   string `json:"credit_token" binding:"required"`
}

type ExternalAPIPayloadItem struct {
	N string      `json:"n"`
	V interface{} `json:"v"`
	T int64       `json:"t"`
}

type ExternalAPIPayload []ExternalAPIPayloadItem

type MQTTInjectionResponse struct {
	N string      `json:"n"`
	V interface{} `json:"v"`
	T int64       `json:"t"`
}

type MQTTInjectionResponsePayload []MQTTInjectionResponse
