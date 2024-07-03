package store

import (
	"errors"
	"os/user"
	"path"

	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/store/sqlite"
	"github.com/spf13/viper"
)

func NewStore() (core.Store, error) {
	switch viper.GetString("database.type") {
	case "sqlite":
		u, _ := user.Current()
		return sqlite.NewStore(path.Join(u.HomeDir, ".pomo/pomo.db"))
	default:
		return nil, errors.New("unknown store type: " + viper.GetString("database.type"))
	}
}
