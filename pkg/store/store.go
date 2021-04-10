package store

import (
	"errors"
	"os/user"
	"path"

	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/joao.rufino/pomo/pkg/store/sqlite"
	"github.com/knadh/koanf"
)

func NewStore(k *koanf.Koanf) (core.Store, error) {
	switch k.String("database.type") {
	case "sqlite":
		u, _ := user.Current()
		return sqlite.NewStore(path.Join(u.HomeDir, ".pomo/pomo.db"))
	default:
		return nil, errors.New("unknown store type" + k.String("database.type"))
	}
}
