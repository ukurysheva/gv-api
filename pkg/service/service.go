package service

import (
	gvapi "github.com/ukurysheva/gv-api"
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

type Authorization interface {
}
type Aircraft interface {
}
type Country interface {
	Create(userId int, country gvapi.Country) (int, error)
	GetAll() ([]gvapi.Country, error)
	GetById(countryId int) (gvapi.Country, error)
	Update(userId, countryId int, info gvapi.Country) error
}
type Airport interface {
}
type Flight interface {
}
type Airline interface {
}
type Ticket interface {
}
