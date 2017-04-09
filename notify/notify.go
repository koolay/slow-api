package notify

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/koolay/slow-api/logging"
)

type ResponseTimeNotification struct {
	Url                  string
	RequestType          string
	ExpectedResponsetime int64
	MeanResponseTime     int64
}

type SlowSqlNotification struct {
	Sql          string
	Host         string  `db:"host"`
	QueryTime    float32 `db:"query_time"`
	LockTime     float32 `db:"lock_time"`
	RowsSent     int32   `db:"rows_sent"`
	RowsExamined int32   `db:"rows_examined"`
}

var (
	errorCount       = 0
	notificationsMap = make(map[string]Notify)
)

type Notify interface {
	GetClientName() string
	Initialize() error
	SendResponseTimeNotification(notification *ResponseTimeNotification) error
	SendSlowSqlNotification(notification *SlowSqlNotification) error
}

func AddNew(name string, notifyItem Notify) {
	notificationsMap[name] = notifyItem
	for _, value := range notificationsMap {
		initErr := value.Initialize()

		if initErr != nil {
			println("Notifications : Failed to Initialize ", value.GetClientName(), ".Please check the details in config file ")
			println("Error Details :", initErr.Error())
		} else {
			println("Notifications :", value.GetClientName(), " Intialized")
		}

	}
}

//Send response time notification to all clients registered
func SendResponseTimeNotification(responseTimeNotification *ResponseTimeNotification) {

	for _, value := range notificationsMap {
		err := value.SendResponseTimeNotification(responseTimeNotification)

		if err != nil {
			logging.Logger.ERROR.Println(err.Error())
		}
	}
}

func SendSlowSqlNotification(slowNotification *SlowSqlNotification) {
	for _, value := range notificationsMap {
		logging.Logger.INFO.Println("send slow mysql mail")
		err := value.SendSlowSqlNotification(slowNotification)

		if err != nil {
			logging.Logger.ERROR.Println(err.Error())

		}
	}
}

//Send Test notification to all registered clients .To make sure everything is working
func SendTestNotification() {

	println("Sending Test notifications to the registered clients")

	for _, value := range notificationsMap {
		err := value.SendResponseTimeNotification(&ResponseTimeNotification{"http://test.com", "GET", 700, 800})

		if err != nil {
			println("Failed to Send Response Time notification to ", value.GetClientName(), " Please check the details entered in the config file")
			println("Error Details :", err.Error())
			os.Exit(3)
		} else {
			println("Sent Test Response Time notification to ", value.GetClientName(), ".Make sure you recieved it")
		}
	}
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

//A readable message string from responseTimeNotification
func getMessageFromResponseTimeNotification(responseTimeNotification *ResponseTimeNotification) string {

	message := fmt.Sprintf("Notification From StatusOk\n\nOne of your apis response time is below than expected."+
		"\n\nPlease find the Details below"+
		"\n\nUrl: %v \nRequestType: %v \nCurrent Average Response Time: %v ms\nExpected Response Time: %v ms\n"+
		"\n\nThanks", responseTimeNotification.Url, responseTimeNotification.RequestType, responseTimeNotification.MeanResponseTime, responseTimeNotification.ExpectedResponsetime)

	return message
}
