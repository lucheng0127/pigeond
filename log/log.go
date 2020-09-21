package log

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Log pigeond logger
var Log = log.New()

// ConfLog configure logger
func ConfLog(fn string, debug bool) error {
	logFile, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	Log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	Log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	Log.SetLevel(log.InfoLevel)
	if debug == true {
		Log.SetLevel(log.DebugLevel)
	}

	return nil
}
