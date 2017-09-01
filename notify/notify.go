package notify

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/koolay/slow-api/logging"
)

var (
	errorCount       = 0
	notificationsMap = make(map[string]Notify)
)

type NotifyItem interface {
	MustNotify() error
}

type SlowAPI struct {
	Url          string
	Method       string
	Responsetime int64
}

func (p *SlowAPI) MustNotify() error {
  return sendSlowAPINotification(p)
}

type SlowSql struct {
	Sql          string
	Host         string    `db:"host"`
	QueryTime    float32   `db:"query_time"`
	LockTime     float32   `db:"lock_time"`
	RowsSent     int32     `db:"rows_sent"`
	RowsExamined int32     `db:"rows_examined"`
	When         time.Time `db:"when"`
}

func (p *SlowSql) MustNotify() error {
  return sendSlowSqlNotification(p)
}

type Notify interface {
	GetClientName() string
	SendSlowAPINotification(notification *SlowAPI) error
	SendSlowSqlNotification(notification *SlowSql) error
}

func AddNew(name string, notifyItem Notify) {
	notificationsMap[name] = notifyItem
}

//Send response time notification to all clients registered
func sendSlowAPINotification(slowAPI *SlowAPI) error {
	var wg sync.WaitGroup
	var lastErr error
	for _, ntf := range notificationsMap {
		wg.Add(1)
		go func(ntf Notify) {
			wg.Done()
			err := ntf.SendSlowAPINotification(slowAPI)
			if err != nil {
				lastErr = err
			}
		}(ntf)
	}
	wg.Wait()
	return lastErr
}

func sendSlowSqlNotification(slowSql *SlowSql) error {
	var wg sync.WaitGroup
	var lastErr error
	for _, ntf := range notificationsMap {
		logging.Logger.INFO.Println("send slow mysql mail")
		wg.Add(1)
		go func(ntf Notify) {
			defer wg.Done()
			err := ntf.SendSlowSqlNotification(slowSql)
			if err != nil {
				lastErr = err
			}
		}(ntf)
	}
	wg.Wait()
	return lastErr
}

//Send Test notification to all registered clients .To make sure everything is working
func SendTestNotification() {
	var wg sync.WaitGroup
	for _, ntf := range notificationsMap {
		wg.Add(2)
		go func(ntf Notify) {
			defer wg.Done()
			err := ntf.SendSlowAPINotification(&SlowAPI{"http://test.com", "GET", 700})
			if err != nil {
				logging.Logger.ERROR.Println(err.Error())
			}
		}(ntf)
		go func(ntf Notify) {
			err := ntf.SendSlowSqlNotification(&SlowSql{Sql: "select test"})
			if err != nil {
				logging.Logger.ERROR.Println(err.Error())
			}
		}(ntf)
	}
	wg.Wait()
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func isEmptyObject(objectString string) bool {
	objectString = strings.Replace(objectString, "0", "", -1)
	objectString = strings.Replace(objectString, "map", "", -1)
	objectString = strings.Replace(objectString, "[]", "", -1)
	objectString = strings.Replace(objectString, " ", "", -1)

	if len(objectString) > 2 {
		return false
	} else {
		return true
	}
}
