package main

import (
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/app"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/config"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/logger/handlers/logruspretty"
	"github.com/natefinch/lumberjack"
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
	// Создаем новый обработчик для записи в файл
	fileHandler := &lumberjack.Logger{
		Filename:   "logs/logger.log",
		MaxSize:    10,   // Максимальный размер файла в мегабайтах
		MaxBackups: 3,    // Максимальное количество ротированных файлов
		MaxAge:     7,    // Максимальный возраст ротированных файлов в днях
		Compress:   true, // Сжатие ротированных файлов
	}

	switch env {
	case envLocal:
		return setupPrettySlog(log)
	case envDev:
		log.SetOutput(fileHandler)
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})

	case envProd:
		log.SetOutput(fileHandler)
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
		log.SetLevel(logrus.WarnLevel)
	default:
		log.SetOutput(fileHandler)
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
		log.SetLevel(logrus.DebugLevel)
	}
	log.SetOutput(os.Stdout)

	return logrus.NewEntry(log)
}

func setupPrettySlog(log *logrus.Logger) *logrus.Entry {
	prettyHandler := logruspretty.NewPrettyHandler(os.Stdout)
	log.SetFormatter(prettyHandler)
	return logrus.NewEntry(log)
}
