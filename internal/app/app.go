package app

import (
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/adapters/db/postgresql"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/adapters/db/redis"
	grpcapp "github.com/markgregr/FruitfulFriends-gRPC-server/internal/app/grpc"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/config"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/services/auth"
	"github.com/markgregr/FruitfulFriends-gRPC-server/pkg/gmiddleware"
	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *logrus.Entry, cfg *config.Config) *App {
	postgre, err := postgresql.New(log.Logger, &cfg.Postgres)
	if err != nil {
		panic(err)
	}

	redis, err := redis.New(log.Logger, &cfg.Redis)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log.Logger, postgre, redis, postgre, postgre, cfg.JWT.TokenTTL)

	authMd := gmiddleware.NewAuthInterceptor(cfg.JWT.TokenKey, authService)

	grpcApp := grpcapp.New(log, authService, authMd, cfg.GRPC.Port)

	return &App{
		GRPCSrv: grpcApp,
	}

}
