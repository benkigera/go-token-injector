package routes

import (
	"github.com/gin-gonic/gin"
	"mqqt_go/api/updater/handlers"
)

func SetupMeterDataRoutes(r *gin.Engine, handler *handlers.MeterDataHandler) {
	r.GET("/api/meter_data", handler.GetMeterData)
	r.OPTIONS("/api/meter_data", handler.GetMeterData) // Handle preflight for CORS
}
