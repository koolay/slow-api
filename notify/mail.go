package notify

import (
	"bytes"
	"encoding/json"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	"github.com/koolay/slow-api/logging"
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
	isAuthorized bool
	client       *smtp.Client
)

func (mailNotify MailNotify) GetClientName() string {
	return "Smtp Mail"
}

func (mailNotify MailNotify) Initialize() error {

	// Check if server listens on that port.
	if len(mailNotify.Username) == 0 && len(mailNotify.Password) == 0 {
		isAuthorized = false

		conn, err := smtp.Dial(mailNotify.Host + ":" + strconv.Itoa(mailNotify.Port))

		if err != nil {
			return err
		}

		client = conn

	} else {
		isAuthorized = true
		conn, err := net.DialTimeout("tcp", mailNotify.Host+":"+strconv.Itoa(mailNotify.Port), 3*time.Second)
		if err != nil {
			return err
		}
		if conn != nil {
			defer conn.Close()
		}
	}
	// Validate sender and recipient
	_, err := mail.ParseAddress(mailNotify.From)
	if err != nil {
		return err
	}

	for _, addr := range mailNotify.To {
		_, err = mail.ParseAddress(addr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mailNotify MailNotify) SendSlowSqlNotification(notification *SlowSqlNotification) error {
	logging.Logger.INFO.Printf("send mail %v", notification)
	message, err := json.Marshal(notification)
	if isAuthorized {

		auth := smtp.PlainAuth("", mailNotify.Username, mailNotify.Password, mailNotify.Host)
		if err == nil {

			// Connect to the server, authenticate, set the sender and recipient,
			// and send the email all in one step.
			err = smtp.SendMail(
				mailNotify.Host+":"+strconv.Itoa(mailNotify.Port),
				auth,
				mailNotify.From,
				mailNotify.To,
				message,
			)
			logging.Logger.FATAL.Fatalln("sent and exit")

			if err != nil {
				return err
			}
			return nil
		}

	} else {
		wc, err := client.Data()

		if err != nil {
			return err
		}

		defer wc.Close()
		if _, err = wc.Write(message); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (mailNotify MailNotify) SendResponseTimeNotification(responseTimeNotification *ResponseTimeNotification) error {
	if isAuthorized {

		auth := smtp.PlainAuth("", mailNotify.Username, mailNotify.Password, mailNotify.Host)
		message := getMessageFromResponseTimeNotification(responseTimeNotification)

		// Connect to the server, authenticate, set the sender and recipient,
		// and send the email all in one step.
		err := smtp.SendMail(
			mailNotify.Host+":"+strconv.Itoa(mailNotify.Port),
			auth,
			mailNotify.From,
			mailNotify.To,
			bytes.NewBufferString(message).Bytes(),
		)

		if err != nil {
			return err
		}
		return nil
	} else {
		wc, err := client.Data()

		if err != nil {
			return err
		}

		defer wc.Close()

		message := bytes.NewBufferString(getMessageFromResponseTimeNotification(responseTimeNotification))

		if _, err = message.WriteTo(wc); err != nil {
			return err
		}

		return nil
	}
}
