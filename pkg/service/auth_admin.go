package service

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

const (
	tokenExp = 12 * time.Hour
)

type AuthAdminService struct {
	repo repository.AuthorizationAdmin
}

func NewAuthAdminService(repo repository.AuthorizationAdmin) *AuthAdminService {
	return &AuthAdminService{repo: repo}
}

func (s *AuthAdminService) CreateAdminUser(adminUser gvapi.AuthAdminUser) (int, error) {
	adminUser.Password = generateAdminPassword(adminUser.Password)
	return s.repo.CreateAdminUser(adminUser)
}

func (s *AuthAdminService) GetUserAdmin(username, password string) (gvapi.AuthAdminUser, error) {
	return s.repo.GetUserAdmin(username, generateAdminPassword(password))
}

func generateAdminPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("ADMIN_SALT"))))
}
