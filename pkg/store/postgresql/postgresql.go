package postgresql

import (
	gormStore "github.com/joaorufino/pomo/pkg/store/gorm"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStore struct {
	*gormStore.GormStore
}

func NewStore(dsn string, logger *zap.SugaredLogger) (*PostgresStore, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to open PostgreSQL database", zap.Error(err))
		return nil, err
	}

	store := gormStore.NewStore(db, logger)

	return &PostgresStore{store}, nil
}

func (s *PostgresStore) InitDB() error {
	return s.GormStore.Init()
}
