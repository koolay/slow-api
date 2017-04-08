package parse

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Parser parse mysql slowlog to go struct.
type MysqlParser struct {
	Location *time.Location
}

type SlowQuery struct {
	ParseTime    time.Time `db:"parse_time"`
	User         string    `db:"user"`
	Host         string    `db:"host"`
	QueryTime    float32   `db:"query_time"`
	LockTime     float32   `db:"lock_time"`
	RowsSent     int32     `db:"rows_sent"`
	RowsExamined int32     `db:"rows_examined"`
	Sql          string    `db:"sql"`
}

var (
	reg1 *regexp.Regexp = regexp.MustCompile(`^#? User\@Host:\s+(\S+)\s+\@\s+(\S+).*`)
	reg2 *regexp.Regexp = regexp.MustCompile(`^# Query_time: ([0-9.]+)\s+Lock_time: ([0-9.]+)\s+Rows_sent: ([0-9.]+)\s+Rows_examined: ([0-9.]+).*`)
)

func (p *SlowQuery) reset() {
	p.Sql = ""
}

// AsLTSV returns parsed slowlog as LTSV format.
func (p *SlowQuery) AsLTSV() string {
	return strings.Join([]string{
		fmt.Sprintf("parse_time:%s", p.ParseTime),
		fmt.Sprintf("user:%s", p.User),
		fmt.Sprintf("host:%s", p.Host),
		fmt.Sprintf("query_time:%f", p.QueryTime),
		fmt.Sprintf("lock_time:%f", p.LockTime),
		fmt.Sprintf("rows_sent:%d", p.RowsSent),
		fmt.Sprintf("rows_examined:%d", p.RowsExamined),
		fmt.Sprintf("sql:%s", p.Sql),
	}, "\t")
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
	if strings.HasPrefix(line, "# Time:") {
		parsed.reset()
		parsed.ParseTime = time.Now() // t.Unix()
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
