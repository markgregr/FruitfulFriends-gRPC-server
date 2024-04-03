package redis

import (
	"github.com/go-redis/redis"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/config"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	rd *redis.Client
}

// New создает новый экземпляр Redis
func New(log *logrus.Logger, cfg *config.RedisConfig) (*Redis, error) {
	const op = "Redis.New"

	log.WithField("op", op).Info("execute redis connection")

	redis := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Password: cfg.Password,
		DB:       0,
	})

	return &Redis{
		rd: redis,
	}, nil
}
