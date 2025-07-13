package repositories

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"mqqt_go/api/updater/models"
)

type MeterDataRepository interface {
	GetLatestMeterData() (*models.MeterData, error)
}

type meterDataRepository struct {
	filePath string
}

func NewMeterDataRepository(filePath string) MeterDataRepository {
	return &meterDataRepository{filePath: filePath}
}

func (r *meterDataRepository) GetLatestMeterData() (*models.MeterData, error) {
	content, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read meter data file: %w", err)
	}

	contentStr := string(content)

	// Regex to extract timestamp
	timestampRegex := regexp.MustCompile(`--- Update at (.*?) ---`)
	timestampMatch := timestampRegex.FindStringSubmatch(contentStr)
	var timestamp time.Time
	if len(timestampMatch) > 1 {
		timestampStr := strings.TrimSpace(timestampMatch[1])
		timestamp, err = time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
	} else {
		// If timestamp not found, use current time or return error
		timestamp = time.Now()
	}

	// Regex to extract Available Credit and Unit
	creditRegex := regexp.MustCompile(`Available Credit: ([0-9.]+) (.*)`)
	creditMatch := creditRegex.FindStringSubmatch(contentStr)

	if len(creditMatch) < 3 {
		return nil, fmt.Errorf("could not parse available credit from file")
	}

	creditValue, err := strconv.ParseFloat(creditMatch[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse available credit value: %w", err)
	}
	unit := strings.TrimSpace(creditMatch[2])

	return &models.MeterData{
		AvailableCredit: creditValue,
		Unit:            unit,
		Timestamp:       timestamp,
	}, nil
}
