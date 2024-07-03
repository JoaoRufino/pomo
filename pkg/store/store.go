package store

import (
	"errors"
	"os/user"
	"path"

	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/store/sqlite"
	"github.com/spf13/viper"
)

func NewStore(conf *conf.Config) (core.Store, error) {
	switch conf.Database.Type {
	case "sqlite":
		u, _ := user.Current()
		return sqlite.NewStore(path.Join(u.HomeDir, conf.Database.Path))
	default:
		return nil, errors.New("unknown store type: " + viper.GetString("database.type"))
	}
}
