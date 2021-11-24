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
			admin := auth.Group("/admin")
			{
				admin.POST("/sign-up", h.adminSignUp)
				admin.POST("/sign-in", h.adminSignIn)
				admin.POST("/sign-out", h.adminLogout)
				admin.POST("/token/refresh", h.tokenAdminRefresh)
			}

			user := auth.Group("/user")
			{
				user.POST("/sign-up", h.userSignUp)
				user.POST("/sign-in", h.userSignIn)
				user.POST("/sign-out", h.userLogout)
				user.POST("/token/refresh", h.tokenUserRefresh)
			}

		}
		countries := api.Group("/countries")
		{
			countries.GET("/", h.getAllCountries)
			countries.GET("/:id", h.getCountryById)

			authenticated := countries.Group("/", h.adminIdentify)
			{
				authenticated.POST("/", h.createCountry)
			}

		}
		airports := api.Group("/airports")
		{
			airports.GET("/", h.getAllAirports)
			airports.GET("/:id", h.getAirportById)

			authenticated := airports.Group("/", h.adminIdentify)
			{
				authenticated.POST("/", h.createAirport)
			}

		}

		airlines := api.Group("/airlines")
		{
			airlines.GET("/", h.getAllAirlines)
			airlines.GET("/:id", h.getAirlineById)

			authenticated := airlines.Group("/", h.adminIdentify)
			{
				authenticated.POST("/", h.createAirline)
			}

		}

		aircrafts := api.Group("/aircrafts")
		{
			aircrafts.GET("/", h.getAllAircrafts)
			aircrafts.GET("/:id", h.getAircraftById)

			authenticated := aircrafts.Group("/", h.adminIdentify)
			{
				authenticated.POST("/", h.createAircraft)
			}

		}

		flights := api.Group("/flights")
		{
			flights.GET("/", h.getAllFlights)
			flights.GET("/:id", h.getFlightById)

			authenticated := flights.Group("/", h.adminIdentify)
			{
				authenticated.POST("/", h.createFlight)
			}

			flights.POST("/search", h.getFlightsByParams)
		}

		users := api.Group("/users")
		{
			authenticated := users.Group("/", h.userIdentify)
			{
				authenticated.POST("/", h.updateUser)
				authenticated.GET("/", h.getUserProfile)

				purchases := authenticated.Group("/purchases")
				{
					purchases.POST("/", h.createPurchase)
					purchases.GET("/:id", h.getPurchaseById)
				}
			}

		}

	}

	return router
}
