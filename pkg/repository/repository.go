package repository

import (
	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type Repository struct {
	AuthorizationAdmin
	User
	Aircraft
	Country
	Airport
	Flight
	Airline
	Purchase
}
type AuthorizationAdmin interface {
	CreateAdminUser(gvapi.AuthAdminUser) (int, error)
	GetUserAdmin(username, password string) (gvapi.AuthAdminUser, error)
}

type User interface {
	CreateUser(gvapi.User) (int, error)
	GetUser(username, password string) (gvapi.User, error)
	Update(userId int, input gvapi.UpdateUserInput) error
	GetProfile(userId int) (gvapi.User, error)
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
	GetByCountryId(countryId int) ([]gvapi.Airport, error)
}

type Airline interface {
	Create(userId int, airline gvapi.Airline) (int, error)
	GetAll() ([]gvapi.Airline, error)
	GetById(airlineId int) (gvapi.Airline, error)
}

type Purchase interface {
	Create(userId int, purchase gvapi.Purchase, flight gvapi.Flight) (int, error)
	GetById(purchaseId int) (gvapi.Purchase, error)
	Update(input gvapi.PurchasePayInput) error
	GetByUserId(userId int) ([]gvapi.Purchase, error)
	GetBasketByUserId(userId int) ([]gvapi.Purchase, error)
}

type Flight interface {
	Create(userId int, airline gvapi.Flight) (int, error)
	GetAll() ([]gvapi.Flight, error)
	GetById(flightId int) (gvapi.Flight, error)
	GetByParams(flightParams gvapi.FlightSearchParams) ([]gvapi.Flight, []gvapi.Flight, error)
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthorizationAdmin: NewAuthAdminPostgres(db),
		User:               NewUserPostgres(db),
		Aircraft:           NewAircraftPostgres(db),
		Airport:            NewAirportPostgres(db),
		Airline:            NewAirlinePostgres(db),
		Country:            NewCountryPostgres(db),
		Flight:             NewFlightPostgres(db),
		Purchase:           NewPurchasePostgres(db),
	}
}
