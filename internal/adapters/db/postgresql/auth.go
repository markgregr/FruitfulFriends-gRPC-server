package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/domain/models"
	"gorm.io/gorm"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

var (
	RoleUser  = 1
	RoleAdmin = 2
)

var (
	StatusActive  = 1
	StatusDeleted = 2
)

func (p *Postgres) SaveUser(ctx context.Context, email string, passHash []byte) (userID int64, err error) {
	const op = "postgresql.Postgres.SaveUser"

	var user = &models.User{
		Email:    email,
		PassHash: passHash,
		Status:   StatusActive,
	}

	if err := p.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID = user.ID

	return userID, nil
}

func (p *Postgres) User(ctx context.Context, email string) (models.User, error) {
	const op = "postgresql.Postgres.User"

	var user models.User

	if err := p.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (p *Postgres) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "postgresql.Postgres.IsAdmin"

	var user models.User

	if err := p.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}
	var isAdmin bool

	switch user.Role {
	case RoleAdmin:
		isAdmin = true
	case RoleUser:
		isAdmin = false
	}

	return isAdmin, nil
}

func (p *Postgres) App(ctx context.Context, appID int) (models.App, error) {
	const op = "postgresql.Postgres.App"

	var app models.App

	if err := p.db.WithContext(ctx).First(&app, appID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app, fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}

		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (p *Postgres) DeleteUser(ctx context.Context, userID int64) error {
	const op = "postgresql.Postgres.DeleteUser"

	var user models.User

	if err := p.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	user.Status = StatusDeleted

	if err := p.db.WithContext(ctx).Save(&user).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Postgres) UpdateUser(ctx context.Context, user *models.User) error {
	const op = "postgresql.Postgres.UpdateUser"

	if err := p.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
