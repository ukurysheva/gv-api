package service

import (
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

type AirportService struct {
	repo        repository.Airport
	countryRepo repository.Country
}

func NewAirportService(repo repository.Airport, countryRepo repository.Country) *AirportService {
	return &AirportService{repo: repo, countryRepo: countryRepo}
}

func (s *AirportService) Create(userId int, airport gvapi.Airport) (int, error) {

	_, err := s.countryRepo.GetById(airport.CountryId)
	if err != nil {
		// country doesn't exist
		return 0, err
	}
	return s.repo.Create(userId, airport)
}

func (s *AirportService) GetAll() ([]gvapi.Airport, error) {
	return s.repo.GetAll()
}

func (s *AirportService) GetById(airport int) (gvapi.Airport, error) {
	return s.repo.GetById(airport)
}

func (s *AirportService) GetByCountryId(CountryId int) ([]gvapi.Airport, error) {
	return s.repo.GetByCountryId(CountryId)
}
