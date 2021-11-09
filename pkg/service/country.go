package service

import (
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

type CountryService struct {
	repo repository.Country
}

func NewCountryService(repo repository.Country) *CountryService {
	return &CountryService{repo: repo}
}

func (s *CountryService) Create(userId int, country gvapi.Country) (int, error) {
	return s.repo.Create(userId, country)
}
func (s *CountryService) GetAll() ([]gvapi.Country, error) {
	return s.repo.GetAll()
}
