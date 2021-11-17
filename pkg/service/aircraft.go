package service

import (
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

type AircraftService struct {
	repo repository.Aircraft
}

func NewAircraftService(repo repository.Aircraft) *AircraftService {
	return &AircraftService{repo: repo}
}

func (s *AircraftService) Create(userId int, country gvapi.Aircraft) (int, error) {
	return s.repo.Create(userId, country)
}
func (s *AircraftService) GetAll() ([]gvapi.Aircraft, error) {
	return s.repo.GetAll()
}

func (s *AircraftService) GetById(countryId int) (gvapi.Aircraft, error) {
	return s.repo.GetById(countryId)
}
