package sqlite

import (
	"os"
	"path/filepath"

	gormStore "github.com/joaorufino/pomo/pkg/store/gorm"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteStore struct {
	*gormStore.GormStore
}

func NewStore(dsn string, logger *zap.SugaredLogger) (*SqliteStore, error) {
	// Ensure directory exists
	dir := filepath.Dir(dsn)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("failed to create directory for SQLite database", zap.Error(err))
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to open SQLite database", zap.Error(err))
		return nil, err
	}

	store := gormStore.NewStore(db, logger)

	return &SqliteStore{store}, nil
}

func (s *SqliteStore) InitDB() error {
	return s.GormStore.Init()
}
