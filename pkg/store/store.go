package store

import (
	"fmt"

	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/store/postgresql"
	"github.com/joaorufino/pomo/pkg/store/sqlite"
	"go.uber.org/zap"
)

func NewStore(conf *conf.Config, logger *zap.SugaredLogger) (core.Store, error) {
	switch conf.Database.Type {
	case "sqlite":
		return sqlite.NewStore(conf.Database.Path, logger)
	case "postgres":
		return postgresql.NewStore(conf.Database.Path, logger)
	default:
		err := fmt.Errorf("unknown store type: %s", conf.Database.Type)
		return nil, err
	}
}
