package grpcapp

import (
	"fmt"
	authgrpc "github.com/markgregr/FruitfulFriends-gRPC-server/internal/grpc/auth"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/logger/sl"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app/grpc.Run"

	log := a.log.With(
		slog.String("op", "app/grpc"),
		slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Error("Failed to listen", sl.Err(err))
		return err
	}

	log.Info("gRPC server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("Failed to serve", sl.Err(err))
		return err
	}
	return nil
}

func (a *App) Stop() {
	const op = "app/grpc.Stop"
	a.log.With(slog.String("op", op)).Info("Stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
