package service

import (
	"errors"
	"forum-project/internal/auth"
	"forum-project/internal/models"
	"forum-project/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password")
var ErrMissmatchPassword = errors.New("passwords do not match")

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (u *UserService) Authenticate(email, password string) (string, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrWrongPassword
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) Register(username, email, password1, password2 string) error {
	if password1 != password2 {
		return ErrMissmatchPassword
	}

	user := models.NewUser()
	user.Username = username
	user.Email = email

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.Role = "user"

	_, err = u.repository.InsertUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := u.repository.GetUserByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
