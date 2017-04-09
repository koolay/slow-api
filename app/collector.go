package app

import (
	"github.com/koolay/slow-api/config"
	"github.com/koolay/slow-api/notify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func getNotifyOptions(notifyName string) (map[string]interface{}, error) {
	notifies := viper.GetStringMap("notifies")
	if notifies != nil {
		if options, ok := notifies[notifyName]; ok {
			return options.(map[string]interface{}), nil
		}
	}
	return nil, errors.Errorf("%s not exist", notifyName)
}

func InitNotification() {
	for _, notifyName := range config.Context.Notify {

		if options, err := getNotifyOptions(notifyName); err == nil {

			switch notifyName {
			case "slack":
				break
			case "mail":
				to := options["to"].([]interface{})
				var toArr []string
				for _, toStr := range to {
					toArr = append(toArr, toStr.(string))
				}
				mailNotify := notify.MailNotify{
					Username: options["username"].(string),
					Password: options["password"].(string),
					Host:     options["host"].(string),
					Port:     int(options["port"].(int64)),
					From:     options["from"].(string),
					To:       toArr,
				}
				notify.AddNew(notifyName, mailNotify)
			}
		}
	}
}
