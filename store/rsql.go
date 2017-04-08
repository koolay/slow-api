package store

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/koolay/slow-api/logging"
	"github.com/koolay/slow-api/parse"
	_ "github.com/lib/pq"
)

var (
	tableName = "slow_sql"
)

type rsqlStorage struct {
	driver string
	dsn    string
}

func newRSqlStorage(driver string, dsn string) (*rsqlStorage, error) {
	if dsn == "" {
		return nil, errors.New("dsn not allow empty")
	}
	rsql := &rsqlStorage{dsn: dsn, driver: driver}
	logging.Logger.INFO.Println("test sql connection")
	if sess, err := rsql.open(); err == nil {
		defer sess.Close()
		if err = sess.Ping(); err != nil {
			return nil, err
		}
		return rsql, nil
	} else {
		return nil, err
	}
}

func (rs *rsqlStorage) open() (*dbr.Session, error) {
	if conn, err := dbr.Open(rs.driver, rs.dsn, nil); err == nil {
		return conn.NewSession(nil), nil
	} else {
		return nil, err
	}
}

func (rs *rsqlStorage) SaveMysqlSlowLog(parsed *parse.SlowQuery) error {
	conn, err := rs.open()
	if err == nil {
		defer conn.Close()
		conn.InsertInto(tableName).Columns("parse_time",
			"user",
			"host",
			"query_time",
			"sql",
			"lock_time",
			"rows_sent",
			"rows_examined").Record(parsed).Exec()
		logging.Logger.INFO.Println("inserted")
	}
	return err
}
