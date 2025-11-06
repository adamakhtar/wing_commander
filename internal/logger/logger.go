package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func SetupLogger(debug bool) (func() error, error) {
	logLevel := log.FatalLevel
	if debug {
		logLevel = log.DebugLevel
	}

	debugFile, err := os.OpenFile("debug.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err == nil {
		log.SetOutput(debugFile)
		log.SetTimeFormat(time.Kitchen)
		log.SetReportCaller(true)
		log.SetLevel(logLevel)
	} else {
		return nil, fmt.Errorf("failed to open debug file: %w", err)
	}

	return debugFile.Close, nil
}