package main

import (
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/app"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/config"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/logger/handlers/logruspretty"
	"github.com/sirupsen/logrus"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.WithField("config", cfg).Info("Application start!")

	application := app.New(log, cfg)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop
	log.Info("Aplication stopping", slog.Any("signal", sign))

	application.GRPCSrv.Stop()

	log.Info("Application stopped!")
}

func setupLogger(env string) *logrus.Entry {
	var log = logrus.New()
	logFilePath := "logs/logger.log"

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	switch env {
	case envLocal:
		return setupPrettySlog(log)
	case envDev:
		log.SetOutput(logFile)
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true, // Добавляем временные метки к сообщениям
		})
	case envProd:
		log.SetOutput(logFile)
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true, // Добавляем временные метки к сообщениям
		})
		log.SetLevel(logrus.WarnLevel)
	default:
		log.SetOutput(logFile)
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true, // Добавляем временные метки к сообщениям
		})
		log.SetLevel(logrus.DebugLevel)
	}

	return logrus.NewEntry(log)
}

func setupPrettySlog(log *logrus.Logger) *logrus.Entry {
	prettyHandler := logruspretty.NewPrettyHandler(os.Stdout)
	log.SetFormatter(prettyHandler)
	return logrus.NewEntry(log)
}
