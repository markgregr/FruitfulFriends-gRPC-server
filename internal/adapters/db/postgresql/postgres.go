package postgresql

import (
	"fmt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/config"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/domain/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	db *gorm.DB
}

// New создает новый экземпляр Postgres
func New(log *logrus.Logger, cfg *config.PostgresConfig) (*Postgres, error) {
	const op = "Postgres.New"

	l := logger.Default
	//if env != "local"{
	//	l = l.LogMode(logger.Info)
	//}

	log.WithField("op", op).Info("execute database connection")

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

// Migrate выполняет миграции базы данных
func Migrate(log *logrus.Logger, db *gorm.DB) error {
	const op = "Postgres.Migrate"

	log.Info("execute database migrations")

	if err := db.AutoMigrate(&models.User{}, &models.App{}); err != nil {
		log.WithError(err).Error("failed to migrate user model")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user model migrated")

	return nil
}

// TestMigrate выполняет миграции базы данных
func TestMigrate(log *logrus.Logger, db *gorm.DB) error {
	const op = "Postgres.TestMigrate"

	log.Info("execute database migrations")

	if err := db.AutoMigrate(&models.User{}, &models.App{}); err != nil {
		log.WithError(err).Error("failed to migrate user model")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("models migrated")

	app := models.App{
		ID:     1,
		Name:   "test",
		Secret: "test-secret",
	}
	if err := db.Create(&app).Error; err != nil {
		log.WithError(err).Error("failed to insert app data")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("data inserted successfully")

	return nil
}
