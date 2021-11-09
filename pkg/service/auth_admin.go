package service

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

const (
	tokenExp = 12 * time.Hour
)

type token struct {
	jwt.StandardClaims
	AdminId int `json:"admin_id"`
}
type AuthAdminService struct {
	repo repository.AuthorizationAdmin
}

func NewAuthAdminService(repo repository.AuthorizationAdmin) *AuthAdminService {
	return &AuthAdminService{repo: repo}
}

func (s *AuthAdminService) CreateAdminUser(adminUser gvapi.AdminUser) (int, error) {
	adminUser.Password = generatePassword(adminUser.Password)
	return s.repo.CreateAdminUser(adminUser)
}
func (s *AuthAdminService) CreateToken(username, password string) (string, error) {
	admin, err := s.repo.GetUserAdmin(username, generatePassword(password))

	if err != nil {
		return "", err
	}

	admClaims := jwt.MapClaims{}
	admClaims["authorized"] = true
	admClaims["admin_id"] = admin.Id
	admClaims["exp"] = time.Now().Add(tokenExp).Unix()

	adm := jwt.NewWithClaims(jwt.SigningMethodHS256, admClaims)
	token, err := adm.SignedString([]byte(os.Getenv("ADMIN_SIGNING_KEY")))

	if err != nil {
		return "", err
	}
	return token, nil

}

func generatePassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("ADMIN_SALT"))))
}
