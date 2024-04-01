package app

import (
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/adapters/db/postgresql"
	grpcapp "github.com/markgregr/FruitfulFriends-gRPC-server/internal/app/grpc"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/config"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/services/auth"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	postgre, err := postgresql.New(log, &cfg.Postgres)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, postgre, postgre, postgre, cfg.JWT.TokenTTL)

	grpcApp := grpcapp.New(log, authService, cfg.GRPC.Port)

	return &App{
		GRPCSrv: grpcApp,
	}

}
