package service

import (
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

type Service struct {
	Authorization
	Aircraft
	Country
	Airport
	Flight
	Airline
	Ticket
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthAdminService(repo.AuthorizationAdmin),
		Country:       NewCountryService(repo.Country),
		Airport:       NewAirportService(repo.Airport, repo.Country),
		Airline:       NewAirlineService(repo.Airline, repo.Country),
		Aircraft:      NewAircraftService(repo.Aircraft),
		Flight:        NewFlightService(repo.Flight, repo.Aircraft, repo.Airport),
	}
}

type Authorization interface {
	GetUserAdmin(username, password string) (gvapi.AdminUser, error)
	CreateAdminUser(gvapi.AdminUser) (int, error)
}
type Aircraft interface {
	Create(userId int, aircraft gvapi.Aircraft) (int, error)
	GetAll() ([]gvapi.Aircraft, error)
	GetById(aircraftId int) (gvapi.Aircraft, error)
}
type Country interface {
	Create(userId int, country gvapi.Country) (int, error)
	GetAll() ([]gvapi.Country, error)
	GetById(countryId int) (gvapi.Country, error)
	// Update(userId, countryId int, info gvapi.Country) error
}
type Airport interface {
	Create(userId int, airport gvapi.Airport) (int, error)
	GetAll() ([]gvapi.Airport, error)
	GetById(airportId int) (gvapi.Airport, error)
}
type Flight interface {
	Create(userId int, airline gvapi.Flight) (int, error)
	GetAll() ([]gvapi.Flight, error)
	GetById(flightId int) (gvapi.Flight, error)
}
type Airline interface {
	Create(userId int, airport gvapi.Airline) (int, error)
	GetAll() ([]gvapi.Airline, error)
	GetById(airlineId int) (gvapi.Airline, error)
}
type Ticket interface {
}
