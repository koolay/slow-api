package store

import (
	"github.com/koolay/slow-api/config"
	"github.com/koolay/slow-api/logging"
	"github.com/koolay/slow-api/parse"
)

type Storage interface {
	SaveMysqlSlowLog(parsed *parse.SlowQuery) error
}

func NewStorage(backend string) (Storage, error) {

	logging.Logger.INFO.Println("use backend ", backend)
	optionsMap, err := config.GetBackends(backend)
	if err != nil {
		return nil, err
	}

	switch backend {
	case "mysql", "postgres":
		dsn := config.ValueOfMap("dsn", optionsMap, "")
		return newRSqlStorage(backend, dsn)
	default:
		dsn := config.ValueOfMap("dsn", optionsMap, "")
		return newRSqlStorage("mysql", dsn)
	}
}
