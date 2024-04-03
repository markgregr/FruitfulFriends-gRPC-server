package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/domain/models"
	"gorm.io/gorm"
)

var (
	ErrUsersNotFound = errors.New("users not found")
)

func (p *Postgres) UserList(ctx context.Context) ([]models.User, error) {
	const op = "postgresql.Postgres.UserList"

	var users []models.User

	if err := p.db.WithContext(ctx).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUsersNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (p *Postgres) UserByID(ctx context.Context, userID int64) (models.User, error) {
	const op = "postgresql.Postgres.UserByID"

	var user models.User

	if err := p.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
