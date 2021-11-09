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
}
type Country interface {
	Create(userId int, country gvapi.Country) (int, error)
	GetAll() ([]gvapi.Country, error)
	// GetById(countryId int) (gvapi.Country, error)
	// Update(userId, countryId int, info gvapi.Country) error
}
type Airport interface {
}

type Airline interface {
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
