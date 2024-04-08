package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/adapters/db/postgresql"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/adapters/db/redis"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/domain/models"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type AuthService struct {
	log           *logrus.Logger
	userSaver     UserSaver
	authUserSaver AuthenticatedUserSaver
	userProvider  UserProvider
	appProvider   AppProvider
	tokenTTl      time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type AuthenticatedUserSaver interface {
	SaveAuthenticatedUser(ctx context.Context, userID int64, token string) error
	UserIDByToken(ctx context.Context, token string) (int64, error)
	DeleteAuthenticatedUser(ctx context.Context, token string) error
}

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExist          = errors.New("user already exists")
)

// New создает новый экземпляр сервиса авторизации
func New(log *logrus.Logger, userSaver UserSaver, authUserSaver AuthenticatedUserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTl time.Duration) *AuthService {
	return &AuthService{
		log:           log,
		userSaver:     userSaver,
		authUserSaver: authUserSaver,
		userProvider:  userProvider,
		appProvider:   appProvider,
		tokenTTl:      tokenTTl,
	}
}

// Login выполняет аутентификацию пользователя
func (s *AuthService) Login(ctx context.Context, email, password string, appID int) (token string, err error) {
	const op = "auth.Auth.Login"
	log := s.log.WithField("op", op).WithField("email", email)

	log.Info("login user")

	user, err := s.userProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			log.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.WithError(err).Error("failed to get user")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid password", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := s.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, postgresql.ErrAppNotFound) {
			log.Warn("app not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		log.WithError(err).Error("failed to get app")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err = jwt.NewToken(user, app, s.tokenTTl)
	if err != nil {
		log.WithError(err).Error("failed to create token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := s.authUserSaver.SaveAuthenticatedUser(ctx, user.ID, token); err != nil {
		log.WithError(err).Error("failed to save authenticated user")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser регистрирует нового пользователя
func (s *AuthService) RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error) {
	op := "auth.Auth.RegisterNewUser"
	log := s.log.WithField("op", op).WithField("email", email)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("failed to hash password")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err = s.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserExists) {
			log.Warn("user already exists", err)
			return 0, fmt.Errorf("%s: %w", op, ErrUserExist)
		}

		log.WithError(err).Error("failed to save user")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

// IsAdmin проверяет, является ли пользователь администратором
func (s *AuthService) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	const op = "auth.Auth.IsAdmin"

	log := s.log.WithField("op", op).WithField("email", userID)

	log.Info("checking if user is admin")

	isAdmin, err = s.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, postgresql.ErrAppNotFound) {
			log.Warn("user not found", err)
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		log.WithError(err).Error("failed to check if user is admin")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

// AuthByToken выполняет аутентификацию пользователя по токену
func (s *AuthService) AuthByToken(ctx context.Context, token string) (userID int64, err error) {
	const op = "auth.Auth.AuthByToken"
	log := s.log.WithField("op", op).WithField("token", token)

	log.Info("auth by token")

	userID, err = s.authUserSaver.UserIDByToken(ctx, token)
	if err != nil {
		if errors.Is(err, redis.ErrTokenNotFound) {
			log.Warn("token not found", err)
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.WithError(err).Error("failed to get user by token")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

// Logout выполняет выход пользователя
func (s *AuthService) Logout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "auth.Auth.Logout"

	//TODO: token :=
	token := "sdfsdf"

	log := s.log.WithField("op", op).WithField("token", token)

	log.Info("logout")

	if err := s.authUserSaver.DeleteAuthenticatedUser(ctx, token); err != nil {
		log.WithError(err).Error("failed to delete authenticated user")
		return &emptypb.Empty{}, fmt.Errorf("%s: %w", op, err)
	}

	return nil, nil
}
