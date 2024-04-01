package postgresql

import (
	"fmt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/config"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/domain/models"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/logger/sl"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
)

type Postgres struct {
	db *gorm.DB
}

// New создает новый экземпляр Postgres
func New(log *slog.Logger, cfg *config.PostgresConfig) (*Postgres, error) {
	const op = "Postgres.New"

	l := logger.Default
	//if env != "local"{
	//	l = l.LogMode(logger.Info)
	//}

	log.With(slog.String("op", op)).Info("execute database connection")

	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{Logger: l})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if cfg.AutoMigrate {
		log.Info("execute database migrations")
		if err := Migrate(log, db); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &Postgres{
		db: db,
	}, nil
}

func Migrate(log *slog.Logger, db *gorm.DB) error {
	const op = "Postgres.Migrate"

	log.Info("execute database migrations")

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Error("failed to migrate user model", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user model migrated")

	return nil
}
