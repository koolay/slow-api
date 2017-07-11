package app

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/hpcloud/tail"
	"github.com/koolay/slow-api/config"
	"github.com/koolay/slow-api/logging"
	"github.com/koolay/slow-api/notify"
	"github.com/koolay/slow-api/parse"
	"github.com/koolay/slow-api/store"
	"github.com/pkg/errors"
)

type MysqlCollector struct {
	filepath         string
	alertSlowSeconds float32
}

func NewMySqlCollector(cfg *config.Config) *MysqlCollector {
	return &MysqlCollector{filepath: cfg.MysqlLogPath, alertSlowSeconds: cfg.MysqlSlowAlertSeconds}
}

func (collector *MysqlCollector) ImportLogFile(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return errors.Errorf("mysql slow log file:[%s] not exits ", filepath)
	}
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)

	parser := parse.NewParser()
	parsed := &parse.SlowQuery{}
	storage, err := store.NewStorage(config.Context.Backend)
	if err == nil {
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				return err
			}
			if parser.Parse(parsed, string(line)) {
				storage.SaveMysqlSlowLog(parsed)
				collector.notify(parsed)
			}
		}
	}
	return nil
}

func (collector *MysqlCollector) Start() error {
	if _, err := os.Stat(collector.filepath); os.IsNotExist(err) {
		return errors.Errorf("mysql slow log file:[%s] not exits ", collector.filepath)
	}

	seekinfo := tail.SeekInfo{Whence: io.SeekEnd}
	cfg := tail.Config{Follow: true, ReOpen: true, Logger: tail.DiscardingLogger, Location: &seekinfo}
	t, err := tail.TailFile(collector.filepath, cfg)
	if err != nil {
		logging.Logger.ERROR.Println(err.Error())
		return errors.Wrapf(err, "tail file: %s", collector.filepath)
	}

	parser := parse.NewParser()
	parsed := &parse.SlowQuery{}
	storage, err := store.NewStorage(config.Context.Backend)
	if err == nil {
		for line := range t.Lines {
			logging.Logger.DEBUG.Println(line)
			if parser.Parse(parsed, line.Text) {
				storage.SaveMysqlSlowLog(parsed)
				collector.notify(parsed)
			}
		}
	}
	return err
}

func (collector *MysqlCollector) notify(parsed *parse.SlowQuery) {
	if parsed.QueryTime >= collector.alertSlowSeconds {
		logging.Logger.DEBUG.Println("notify")
		body := &notify.SlowSql{
			Sql:          parsed.Sql,
			Host:         parsed.Host,
			QueryTime:    parsed.QueryTime,
			LockTime:     parsed.LockTime,
			RowsSent:     parsed.RowsSent,
			RowsExamined: parsed.RowsExamined,
			When:         parsed.When,
		}
		notify.SendSlowSqlNotification(body)
	}
}
