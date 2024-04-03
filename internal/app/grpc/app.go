package grpcapp

import (
	"fmt"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	authgrpc "github.com/markgregr/FruitfulFriends-gRPC-server/internal/grpc/auth"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/logger/handlers/logruspretty"
	"github.com/markgregr/FruitfulFriends-gRPC-server/pkg/gmiddleware"
	"github.com/markgregr/FruitfulFriends-gRPC-server/pkg/gserver"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

type App struct {
	log        *logrus.Entry
	gRPCServer *grpc.Server
	port       int
}

func New(log *logrus.Entry, authService authgrpc.AuthService, authMd *gmiddleware.Auth, port int) *App { // Создаем экземпляр PrettyHandler для вывода красивых логов
	prettyHandler := logruspretty.NewPrettyHandler(os.Stdout)
	logrus.SetFormatter(prettyHandler)
	logEntry := logrus.NewEntry(logrus.StandardLogger())

	gRPCServer := grpc.NewServer(
		gserver.StdUnaryMiddleware(logEntry, grpcauth.UnaryServerInterceptor(authMd.AuthFunc)),
		gserver.StdStreamMiddleware(logEntry, grpcauth.StreamServerInterceptor(authMd.AuthFunc)),
	)

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

	log := a.log.Logger.WithField("op", op).WithField("port", a.port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.WithError(err).Error("Failed to listen")
		return err
	}

	log.Info("gRPC server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		log.WithError(err).Error("Failed to serve")
		return err
	}
	return nil
}

func (a *App) Stop() {
	const op = "app/grpc.Stop"
	a.log.Logger.WithField("op", op).Info("Stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
