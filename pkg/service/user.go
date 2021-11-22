package service

import (
	"crypto/sha1"
	"fmt"
	"os"

	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user gvapi.User) (int, error) {
	user.Password = generateUserPassword(user.Password)
	id, err := s.repo.CreateUser(user)

	if err != nil {
		return 0, err
	}

	subject := "Subject: Регистрация в GlobalAvia\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := "<html><head><title>Здравствуйте, " + user.FirstName + " " + user.LastName + "!</title></head>" +
		"<body><br>Спасибо за регистрацию в GlobalAvia.<br>Ваши данные для входа: <br><b>Email</b>: " + user.Email + "</body></html>" + " \r\n"

	SendMail(user.Email, subject, mime, msg)
	return id, nil
}

func (s *UserService) GetUser(username, password string) (gvapi.User, error) {
	return s.repo.GetUser(username, generateUserPassword(password))
}
func (s *UserService) GetProfile(userId int) (gvapi.User, error) {
	return s.repo.GetProfile(userId)
}

func (s *UserService) Update(userId int, input gvapi.UpdateUserInput) error {
	return s.repo.Update(userId, input)
}

func generateUserPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("USER_SALT"))))
}
