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

// NewAuth returns auth usecase
func NewAuth(userRepo domain.UserRepository, logger logger.Logger) Auth {
	return &auth{userRepo: userRepo, logger: logger}
}

// SignUp the user
func (a *auth) SignUp(ctx context.Context, user *domain.User) (*domain.User, error) {
	// check if user with given email exists
	if found, err := a.userRepo.FindByEmail(ctx, user.Email); found != nil {
		a.logger.Info(found)
		a.logger.Info(err)
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
	createdUser.HidePassword()

	return createdUser, nil
}

// SignIn the user
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

	// hide user's password in response
	foundUser.HidePassword()

	return foundUser, nil
}

// FindUserByEmail finds the user by their email
func (a *auth) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// check if user with given email exists
	foundUser, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		a.logger.Error(err)
		return nil, rest.NewBadRequestError(rest.ErrWrongCredentials.Error())
	}

	// hide user's password in response
	foundUser.HidePassword()

	return foundUser, nil
}

// Auth contract
type Auth interface {
	SignUp(ctx context.Context, user *domain.User) (*domain.User, error)
	SignIn(ctx context.Context, user *domain.User) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
