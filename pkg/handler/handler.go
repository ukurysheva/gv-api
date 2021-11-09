package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ukurysheva/gv-api/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/v1")
	{
		countries := api.Group("/countries")
		{
			countries.GET("/", h.getAllCountries)
			countries.POST("/", h.userIdentity, h.createCountry)
			// models.GET("/shop/:shop", h.getShopModels)
			// models.GET("/search", h.searchModels)

			// lists.GET("/:id", h.getListById)

		}
		// api.GET("/connect")
	}

	return router
}
