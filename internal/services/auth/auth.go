package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/adapters/db/postgresql"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/domain/models"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/jwt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/logger/sl"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTl     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
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
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTl time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTl:     tokenTTl,
	}
}

// Login выполняет аутентификацию пользователя
func (a *Auth) Login(ctx context.Context, email, password string, appID int) (token string, err error) {
	const op = "auth.Auth.Login"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid password", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, postgresql.ErrAppNotFound) {
			log.Warn("app not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		log.Error("failed to get app", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err = jwt.NewToken(user, app, a.tokenTTl)
	if err != nil {
		log.Error("failed to create token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser регистрирует нового пользователя
func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error) {
	op := "auth.Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err = a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExist)
		}

		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

// IsAdmin проверяет, является ли пользователь администратором
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	const op = "auth.Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err = a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, postgresql.ErrAppNotFound) {
			log.Warn("user not found", sl.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		log.Error("failed to check if user is admin", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
