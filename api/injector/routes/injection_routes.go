package routes

import (
	"github.com/gin-gonic/gin"
	"mqqt_go/api/injector/handlers"
)

func SetupInjectionRoutes(r *gin.Engine, handler *handlers.InjectionHandler) {
	r.POST("/api/inject_token", handler.InjectToken)
}
