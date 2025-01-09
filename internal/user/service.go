package user

import (
	"context"
	"crowdfunding/config"
	"errors"
	"time"
)

type Service interface {
	GetUserByID(ctx context.Context, Id int) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) service {
	return service{
		repo: repo,
	}
}

func (s service) RegisterUser(ctx context.Context, request RegisterUserRequest) (user User, err error) {

	user = NewFromRegisterRequest(request)

	if err = user.EncryptPassword(int(config.Cfg.App.Encryption.Salt)); err != nil {
		return
	}

	newUser, err := s.repo.Save(ctx, user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil

}

func (s service) LoginUser(ctx context.Context, request LoginUserRequest) (user User, err error) {

	user, err = s.repo.FindByEmail(ctx, request.Email)
	if err != nil {
		return
	}
	if err = user.VerifyPassword(request.Password); err != nil {
		err = errors.New("password not match")
		return
	}

	return

}

func (s service) IsEmailAvailable(ctx context.Context, input CheckEmailRequest) (bool, error) {
	email := input.Email

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if user.Id == 0 {
		return true, nil
	}

	return false, nil
}

func (s service) SaveAvatar(ctx context.Context, Id int, fileLocation string) (user User, err error) {

	user, err = s.repo.FindByID(ctx, Id)
	if err != nil {
		return
	}

	user.AvatarFileName = fileLocation
	user.UpdatedAt = time.Now()

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil

	// user, err := s.repo.FindByID(Id)
	// if err != nil {
	// 	return user, err
	// }

	// user.AvatarFileName = fileLocation

	// updatedUser, err := s.repo.Update(user)
	// if err != nil {
	// 	return updatedUser, err
	// }

	// return updatedUser, nil

}

func (s service) GetUserByID(ctx context.Context, Id int) (User, error) {

	user, err := s.repo.FindByID(ctx, Id)
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		return user, errors.New("user not found")
	}

	return user, nil
}
