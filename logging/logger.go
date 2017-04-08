package logging

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/koolay/slow-api/config"
	jww "github.com/spf13/jwalterweatherman"
)

var (
	Logger    *jww.Notepad
	logHandle = os.Stdout
)

func NewLogger(gf *config.Config) *jww.Notepad {
	var logger *jww.Notepad

	var err error
	if gf.LogFile != "" {
		logHandle, err = os.OpenFile(gf.LogFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
	}

	if gf.LogLevel != "" {
		switch strings.ToLower(gf.LogLevel) {
		case "info":
			logger = jww.NewNotepad(jww.LevelInfo, jww.LevelTrace, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
		case "debug":
			logger = jww.NewNotepad(jww.LevelDebug, jww.LevelDebug, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
		case "warn":
			logger = jww.NewNotepad(jww.LevelWarn, jww.LevelTrace, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
		case "error":
			logger = jww.NewNotepad(jww.LevelError, jww.LevelTrace, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
		case "trace":
			logger = jww.NewNotepad(jww.LevelTrace, jww.LevelTrace, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
		default:
			logger = jww.NewNotepad(jww.LevelError, jww.LevelTrace, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
		}
	} else if config.GlobalFlag.Verbose {
		logger = jww.NewNotepad(jww.LevelDebug, jww.LevelDebug, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
	} else {
		logger = jww.NewNotepad(jww.LevelError, jww.LevelTrace, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
	}

	return logger
}
