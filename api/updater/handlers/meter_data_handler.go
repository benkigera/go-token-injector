package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"mqqt_go/api/updater/services"
)

type MeterDataHandler struct {
	service services.MeterDataService
}

func NewMeterDataHandler(service services.MeterDataService) *MeterDataHandler {
	return &MeterDataHandler{
		service: service,
	}
}

// GetMeterData handles the GET request for meter data.
func (h *MeterDataHandler) GetMeterData(c *gin.Context) {
	// Set CORS headers to allow requests from any origin
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type")

	// Handle OPTIONS preflight request
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusOK)
		return
	}

	// Ensure it's a GET request
	if c.Request.Method != "GET" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	meterData, err := h.service.GetLatestMeterData(c.Request.Context())
	if err != nil {
		log.Printf("Error getting meter data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve meter data"})
		return
	}

	c.JSON(http.StatusOK, meterData)
}
