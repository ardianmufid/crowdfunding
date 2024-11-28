package user

import (
	"crowdfunding/config"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) service {
	return service{
		repo: repo,
	}
}

func (s service) RegisterUser(request RegisterUserRequest) (User, error) {

	user := User{}
	user.Name = request.Name
	user.Email = request.Email
	user.Occupation = request.Occupation

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), int(config.Cfg.App.Encryption.Salt))
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repo.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s service) LoginUser(request LoginUserRequest) (User, error) {

	email := request.Email
	password := request.Password

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		return user, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}
