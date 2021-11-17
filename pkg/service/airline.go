package service

import (
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

type AirlineService struct {
	repo        repository.Airline
	countryRepo repository.Country
}

func NewAirlineService(repo repository.Airline, countryRepo repository.Country) *AirlineService {
	return &AirlineService{repo: repo, countryRepo: countryRepo}
}

func (s *AirlineService) Create(userId int, airline gvapi.Airline) (int, error) {

	_, err := s.countryRepo.GetById(airline.CountryId)
	if err != nil {
		// country doesn't exist
		return 0, err
	}
	return s.repo.Create(userId, airline)
}

func (s *AirlineService) GetAll() ([]gvapi.Airline, error) {
	return s.repo.GetAll()
}

func (s *AirlineService) GetById(airlineId int) (gvapi.Airline, error) {
	return s.repo.GetById(airlineId)
}
