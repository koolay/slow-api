package parse

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/koolay/slow-api/logging"
)

// Parser parse mysql slowlog to go struct.
type MysqlParser struct {
	Location *time.Location
}

type SlowQuery struct {
	When         time.Time `db:"when"`
	User         string    `db:"user"`
	Host         string    `db:"host"`
	QueryTime    float32   `db:"query_time"`
	LockTime     float32   `db:"lock_time"`
	RowsSent     int32     `db:"rows_sent"`
	RowsExamined int32     `db:"rows_examined"`
	Sql          string    `db:"sql"`
}

var (
	regTime *regexp.Regexp = regexp.MustCompile(`^#? Time:\s+(\S+)\s+(\S+).*`)
	reg1    *regexp.Regexp = regexp.MustCompile(`^#? User\@Host:\s+(\S+)\s+\@\s+(\S+).*`)
	reg2    *regexp.Regexp = regexp.MustCompile(`^# Query_time: ([0-9.]+)\s+Lock_time: ([0-9.]+)\s+Rows_sent: ([0-9.]+)\s+Rows_examined: ([0-9.]+).*`)
)

func (p *SlowQuery) reset() {
	p.Sql = ""
}

// AsJSON returns parsed slowlog as JSON format.
func (p *SlowQuery) AsJSON() string {
	j, _ := json.Marshal(p)
	return string(j)
}

// NewParser returns new Parser
func NewParser() MysqlParser {
	loc, _ := time.LoadLocation("UTC")
	return MysqlParser{loc}
}

// Parse mysql slowlog.
// You can receive parsed slowlog through channnel.
func (p *MysqlParser) Parse(parsed *SlowQuery, line string) (completed bool) {

	completed = false
	if shouldIgnore(line) {
		return false
	}

	// DateTime
	if r := regTime.FindStringSubmatch(line); r != nil {
		ym := r[1]
		hms := r[2]
		y := "20" + ym[:2]
		m := ym[2:4]
		d := ym[4:6]
		if when, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%s %s", y, m, d, hms)); err == nil {
			parsed.When = when
		} else {
			logging.Logger.ERROR.Println(err.Error())
		}
		parsed.reset()
		return false
	}

	// User, Host
	if r := reg1.FindStringSubmatch(line); r != nil {
		parsed.User = r[1]
		parsed.Host = r[2]
		return false
	}

	// QueryTime, LockTime, RowsSent, RowsExamined
	if r := reg2.FindStringSubmatch(line); r != nil {
		parsed.QueryTime = stringToFloat32(r[1])
		parsed.LockTime = stringToFloat32(r[2])
		parsed.RowsSent = stringToInt32(r[3])
		parsed.RowsExamined = stringToInt32(r[4])
		return false
	}

	// Sql
	if !strings.HasPrefix(line, "#") {
		parsed.Sql += strings.Trim(line, " \r\n") + " "

		if strings.HasSuffix(line, ";") && parsed.Sql != "" {
			parsed.Sql = strings.Trim(parsed.Sql, " ")
			return true
		}
	}
	return false
}

func shouldIgnore(line string) bool {
	if strings.TrimSpace(line) == "" {
		return true
	}
	uppered := strings.ToUpper(line)
	return strings.HasPrefix(uppered, "USE") ||
		strings.HasPrefix(uppered, "SET TIMESTAMP=") ||
		strings.HasSuffix(uppered, "STARTED WITH:") ||
		strings.HasPrefix(uppered, "TIME") ||
		strings.HasPrefix(uppered, "TCP PORT")
}
