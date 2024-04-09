package main

import (
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/adapters/db/postgresql"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/config"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/logger/handlers/logruspretty"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	const op = "main"

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Application start!", slog.Any("config", cfg))

	db, err := gorm.Open(postgres.Open(cfg.Postgres.URL))
	if err != nil {
		panic(err)
	}

	if err = postgresql.TestMigrate(log.Logger, db); err != nil {
		panic(err)
	}
}

func setupLogger(env string) *logrus.Entry {
	var log = logrus.New()
	// Создаем новый обработчик для записи в файл
	fileHandler := &lumberjack.Logger{
		Filename:   "logs/logger_rest.log",
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
