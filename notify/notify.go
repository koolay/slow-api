package notify

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/koolay/slow-api/logging"
)

type SlowAPI struct {
	Url          string
	Method       string
	Responsetime int64
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

var (
	errorCount       = 0
	notificationsMap = make(map[string]Notify)
)

type Notify interface {
	GetClientName() string
	SendSlowAPINotification(notification *SlowAPI) error
	SendSlowSqlNotification(notification *SlowSql) error
}

func AddNew(name string, notifyItem Notify) {
	notificationsMap[name] = notifyItem
}

//Send response time notification to all clients registered
func SendSlowAPINotification(slowAPI *SlowAPI) {
	var wg sync.WaitGroup
	for _, ntf := range notificationsMap {
		wg.Add(1)
		go func(ntf Notify) {
			wg.Done()
			err := ntf.SendSlowAPINotification(slowAPI)
			if err != nil {
				logging.Logger.ERROR.Println(err.Error())
			}
		}(ntf)
	}
	wg.Wait()
}

func SendSlowSqlNotification(slowSql *SlowSql) {
	var wg sync.WaitGroup
	for _, ntf := range notificationsMap {
		logging.Logger.INFO.Println("send slow mysql mail")
		wg.Add(1)
		go func(ntf Notify) {
			defer wg.Done()
			err := ntf.SendSlowSqlNotification(slowSql)
			if err != nil {
				logging.Logger.ERROR.Println(err.Error())
			}
		}(ntf)
	}
	wg.Wait()
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
