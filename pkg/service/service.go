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
	}
}

type Authorization interface {
	CreateToken(username, password string) (string, error)
	CreateAdminUser(gvapi.AdminUser) (int, error)
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
type Flight interface {
}
type Airline interface {
}
type Ticket interface {
}
