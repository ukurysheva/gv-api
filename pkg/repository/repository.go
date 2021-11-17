package repository

import (
	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type Repository struct {
	AuthorizationClient
	AuthorizationAdmin
	Aircraft
	Country
	Airport
	Flight
	Airline
	Ticket
}

type AuthorizationClient interface {
}
type AuthorizationAdmin interface {
	CreateAdminUser(gvapi.AdminUser) (int, error)
	GetUserAdmin(username, password string) (gvapi.AdminUser, error)
}
type Aircraft interface {
	Create(userId int, country gvapi.Aircraft) (int, error)
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

type Airline interface {
	Create(userId int, airline gvapi.Airline) (int, error)
	GetAll() ([]gvapi.Airline, error)
	GetById(airlineId int) (gvapi.Airline, error)
}

type Ticket interface {
}

type Flight interface {
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthorizationClient: NewAuthClientPostgres(db),
		AuthorizationAdmin:  NewAuthAdminPostgres(db),
		Aircraft:            NewAircraftPostgres(db),
		Airport:             NewAirportPostgres(db),
		Airline:             NewAirlinePostgres(db),
		Country:             NewCountryPostgres(db),
	}
}
