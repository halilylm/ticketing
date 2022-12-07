package usecase

import (
	"context"

	"github.com/halilylm/gommon/logger"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/gommon/utils"
	"github.com/halilylm/ticketing/auth/domain"
)

type auth struct {
	userRepo domain.UserRepository
	logger   logger.Logger
}

func (a *auth) SignUp(ctx context.Context, user *domain.User) (*domain.User, error) {
	// check if user with given email exists
	if _, err := a.userRepo.FindByEmail(ctx, user.Email); err != nil {
		a.logger.Error(err)
		return nil, rest.NewBadRequestError(rest.ErrEmailAlreadyExists.Error())
	}

	// hash the password
	user.Password = utils.HashPassword(user.Password)

	// save the user to the database
	createdUser, err := a.userRepo.Insert(ctx, user)

	// return internal server error when user cannot inserted
	if err != nil {
		a.logger.Error(err)
		return nil, rest.NewInternalServerError()
	}

	// hide user's password in response
	createdUser.Password = ""

	return createdUser, nil
}

func (a *auth) SignIn(ctx context.Context, user *domain.User) (*domain.User, error) {
	// check if user with given email exists
	foundUser, err := a.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		a.logger.Error(err)
		return nil, rest.NewBadRequestError(rest.ErrWrongCredentials.Error())
	}

	// compare the password
	if err := utils.CheckPassword(foundUser.Password, user.Password); err != nil {
		a.logger.Error(err)
		return nil, rest.NewBadRequestError(rest.ErrWrongCredentials.Error())
	}

	return foundUser, nil
}

func (a *auth) CurrentUser(ctx context.Context, email string) (*domain.User, error) {
	// check if user with given email exists
	foundUser, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		a.logger.Error(err)
		return nil, rest.NewBadRequestError(rest.ErrWrongCredentials.Error())
	}
	return foundUser, nil
}
