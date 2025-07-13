package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mqqt_go/api/injector/models"
	"mqqt_go/api/injector/services"
)

type InjectionHandler struct {
	service services.InjectionService
}

func NewInjectionHandler(service services.InjectionService) *InjectionHandler {
	return &InjectionHandler{
		service: service,
	}
}

func (h *InjectionHandler) InjectToken(c *gin.Context) {
	var req models.InjectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to handle the token injection logic
	response, err := h.service.InjectToken(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Status == "failed" {
		c.JSON(http.StatusBadRequest, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}
