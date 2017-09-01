package config

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	LogFile               string
	LogLevel              string
	MysqlLogPath          string
	MysqlSlowAlertSeconds float32
	Backend               string
	Notify                []string
	// same slow record, notify duration seconds
	NotifyDuration int
}

var (
	Context           *Config
	supportedBackends map[string]int = map[string]int{"mysql": 1, "postgres": 1, "pg": 1}
)

func InitConfig(cfg *Config) *Config {
	if cfg == nil {
		cfg = &Config{}
		cfg.LogFile = viper.GetString("log_file")
		cfg.LogLevel = viper.GetString("log_level")
		cfg.Backend = viper.GetString("backend")
		cfg.Notify = viper.GetStringSlice("notify")
		cfg.NotifyDuration = viper.GetInt("notify_duration")
	}
	initMySqlCollectorOptions(cfg)
	Context = cfg
	return cfg
}

func initMySqlCollectorOptions(cfg *Config) {
	options := viper.GetStringMapString("collectors.mysql")
	if len(options) < 1 {
		return
	}
	cfg.MysqlLogPath = options["slowlog"]
	if alertSeconds, err := strconv.ParseFloat(options["alert_slow_seconds"], 32); err == nil {
		cfg.MysqlSlowAlertSeconds = float32(alertSeconds)
	} else {
		panic(err)
	}
}

func ValueOfMap(key string, m map[string]interface{}, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

func GetBackends(backend string) (map[string]interface{}, error) {

	backend = strings.ToLower(backend)
	if _, ok := supportedBackends[backend]; !ok {
		return nil, errors.New("not supported backend:" + backend)
	}
	backends := viper.GetStringMap("backends")

	for name, props := range backends {
		if backend == strings.ToLower(name) {
			return props.(map[string]interface{}), nil
		}
	}
	return nil, errors.New("not found configuration of the backend:" + backend)
}
