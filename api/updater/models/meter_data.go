package models

import (
	"time"
)

type MeterData struct {
	AvailableCredit float64   `json:"available_credit"`
	Unit            string    `json:"unit"`
	Timestamp       time.Time `json:"timestamp"`
}
