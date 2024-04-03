package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/adapters/db/postgresql"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/domain/models"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	log          *logrus.Logger
	userProvider UserProvider
}

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UserByID(ctx context.Context, userID int64) (models.User, error)
	UserList(ctx context.Context) ([]models.User, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExist          = errors.New("user already exists")
)

func New(log *logrus.Logger, userProvider UserProvider) *UserService {
	return &UserService{
		log:          log,
		userProvider: userProvider,
	}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "UserService.GetUserByEmail"
	log := s.log.WithField("op", op).WithField("email", email)

	log.Info("getting user by email")

	user, err := s.userProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.WithError(err).Error("failed to get user by email")
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	const op = "UserService.GetUserByID"
	log := s.log.WithField("op", op).WithField("userID", userID)

	log.Info("getting user by ID")

	user, err := s.userProvider.UserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.WithError(err).Error("failed to get user by ID")
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserService) GetUserList(ctx context.Context) ([]models.User, error) {
	const op = "UserService.GetUserList"
	log := s.log.WithField("op", op)

	log.Info("getting user list")

	users, err := s.userProvider.UserList(ctx)
	if err != nil {
		log.WithError(err).Error("failed to get user list", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}
