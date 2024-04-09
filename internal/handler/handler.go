package handler

import (
	"github.com/dkshi/bwgtest/internal/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	quotations := router.Group("/quotations")
	{
		quotations.GET("/update", h.updateQuotation)
		quotations.GET("/get", h.getQuotation)
		quotations.GET("/latest", h.getLatestQuotation)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
