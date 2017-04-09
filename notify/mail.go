package notify

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
)

type MailNotify struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Host     string   `json:"smtpHost"`
	Port     int      `json:"port"`
	From     string   `json:"from"`
	To       []string `json:"to"`
}

var (
	client *smtp.Client
)

func (mailNotify MailNotify) GetClientName() string {
	return "Smtp Mail"
}

func (mailNotify MailNotify) Send(subject, message string) error {

	var content string
	header := make(map[string]string)
	header["From"] = mailNotify.From
	header["To"] = strings.Join(mailNotify.To, ";")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=iso-8859-1"
	for k, v := range header {
		content += fmt.Sprintf("%s: %s\n", k, v)
	}
	content += "<html><body> \n"
	content += message
	content += "</body></html>"

	auth := smtp.PlainAuth("", mailNotify.Username, mailNotify.Password, mailNotify.Host)
	return smtp.SendMail(
		mailNotify.Host+":"+strconv.Itoa(mailNotify.Port),
		auth,
		mailNotify.From,
		mailNotify.To,
		[]byte(content),
	)
}

func (mailNotify MailNotify) SendSlowSqlNotification(notification *SlowSql) error {
	message := fmt.Sprintf(`
	<b>Host:</b> %s </br>
	<b>QueryTime:</b> %f </br>
	<b>LockTime:</b> %f </br>
	<b>RowsSent:</b> %d </br>
	<b>RowsExamined:</b> %d </br>
	<b>When:</b> %s </br>
	<b>Sql:</b> <pre><code style="display:block;white-space:pre-wrap">%s</code></pre> </br>
	`, notification.Host,
		notification.QueryTime,
		notification.LockTime,
		notification.RowsSent,
		notification.RowsExamined,
		notification.When,
		notification.Sql,
	)
	return mailNotify.Send("SlowAPI-slow sql", message)
}

func (mailNotify MailNotify) SendSlowAPINotification(notification *SlowAPI) error {

	message := fmt.Sprintf(`
	<b>URL:</b> %s </br>
	<b>Method:</b> %s </br>
	<b>ResponseTime:</b> %d </br>
	`,
		notification.Url,
		notification.Method,
		notification.Responsetime,
	)

	return mailNotify.Send("SlowAPI-slow api", message)
}
