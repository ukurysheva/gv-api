package service

import (
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

type PurchaseService struct {
	repo        repository.Purchase
	flightsRepo repository.Flight
}

func NewPurchaseService(repo repository.Purchase, flightsRepo repository.Flight) *PurchaseService {
	return &PurchaseService{repo: repo, flightsRepo: flightsRepo}
}

func (s *PurchaseService) Create(userId int, purchase gvapi.Purchase) (int, error) {

	_, err := s.flightsRepo.GetById(purchase.FlightId)
	if err != nil {
		// country doesn't exist
		return 0, err
	}
	return s.repo.Create(userId, purchase)
}

func (s *PurchaseService) GetById(purchaseId int) (gvapi.Purchase, error) {
	return s.repo.GetById(purchaseId)
}

func (s *PurchaseService) GetByUserId(userId int) ([]gvapi.Purchase, error) {
	return s.repo.GetByUserId(userId)
}

func (s *PurchaseService) GetBasketByUserId(userId int) ([]gvapi.Purchase, error) {
	return s.repo.GetBasketByUserId(userId)
}

func (s *PurchaseService) Update(input gvapi.PurchasePayInput) error {
	return s.repo.Update(input)
}
