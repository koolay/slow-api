package store

import (
	"github.com/koolay/slow-api/logging"
	"github.com/koolay/slow-api/parse"
)

func SaveMysqlSlowLog(parsed *parse.SlowQuery) error {
	logging.Logger.INFO.Println(parsed.AsJSON())
	return nil
}
