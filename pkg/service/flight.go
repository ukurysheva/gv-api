package service

import (
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

type FlightService struct {
	repo         repository.Flight
	aircraftRepo repository.Aircraft
	airportRepo  repository.Airport
}

func NewFlightService(repo repository.Flight, aircraftRepo repository.Aircraft, airportRepo repository.Airport) *FlightService {
	return &FlightService{repo: repo, aircraftRepo: aircraftRepo, airportRepo: airportRepo}
}

func (s *FlightService) Create(userId int, airport gvapi.Flight) (int, error) {

	_, err := s.aircraftRepo.GetById(airport.AircraftId)
	if err != nil {
		// aircraft doesn't exist
		return 0, err
	}
	_, err = s.airportRepo.GetById(airport.AirportDepId)
	if err != nil {
		// airport doesn't exist
		return 0, err
	}
	_, err = s.airportRepo.GetById(airport.AirportLandId)
	if err != nil {
		// airport doesn't exist
		return 0, err
	}
	return s.repo.Create(userId, airport)
}

func (s *FlightService) GetAll() ([]gvapi.Flight, error) {
	return s.repo.GetAll()
}

func (s *FlightService) GetById(airport int) (gvapi.Flight, error) {
	flights, err := s.repo.GetById(airport)

	return flights, err
}

func (s *FlightService) GetFlightByParams(input gvapi.FlightSearchParams) ([]gvapi.Flight, []gvapi.Flight, error) {
	return s.repo.GetByParams(input)
}
