package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenInvalid  = errors.New("token invalid")
)

func (r *Redis) SaveAuthenticatedUser(ctx context.Context, userID int64, token string) error {
	return r.rd.Set(token, userID, 0).Err()
}

func (r *Redis) UserIDByToken(ctx context.Context, token string) (int64, error) {
	userID, err := r.rd.Get(token).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, ErrTokenNotFound
		}
		return 0, err
	}

	return userID, nil
}

func (r *Redis) DeleteAuthenticatedUser(ctx context.Context, token string) error {
	return r.rd.Del(token).Err()
}
