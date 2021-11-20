package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/ukurysheva/gv-api/pkg/service"
)

type Handler struct {
	services *service.Service
	redis    *redis.Client
}

func NewHandler(services *service.Service, redis *redis.Client) *Handler {
	return &Handler{services: services, redis: redis}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
			auth.POST("/sign-out", h.userLogout)
			auth.POST("/token/refresh", h.refresh)
		}
		countries := api.Group("/countries")
		{
			countries.GET("/", h.getAllCountries)
			countries.GET("/:id", h.getCountryById)

			authenticated := countries.Group("/", h.userIdentify)
			{
				authenticated.POST("/", h.createCountry)
			}

		}
		airports := api.Group("/airports")
		{
			airports.GET("/", h.getAllAirports)
			airports.GET("/:id", h.getAirportById)

			authenticated := airports.Group("/", h.userIdentify)
			{
				authenticated.POST("/", h.createAirport)
			}

		}

		airlines := api.Group("/airlines")
		{
			airlines.GET("/", h.getAllAirlines)
			airlines.GET("/:id", h.getAirlineById)

			authenticated := airlines.Group("/", h.userIdentify)
			{
				authenticated.POST("/", h.createAirline)
			}

		}

		aircrafts := api.Group("/aircrafts")
		{
			aircrafts.GET("/", h.getAllAircrafts)
			aircrafts.GET("/:id", h.getAircraftById)

			authenticated := aircrafts.Group("/", h.userIdentify)
			{
				authenticated.POST("/", h.createAircraft)
			}

		}

		flights := api.Group("/flights")
		{
			flights.GET("/", h.getAllFlights)
			flights.GET("/:id", h.getFlightById)

			authenticated := flights.Group("/", h.userIdentify)
			{
				authenticated.POST("/", h.createFlight)
			}

		}
	}

	return router
}
