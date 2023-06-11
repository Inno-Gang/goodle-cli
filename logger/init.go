package logger

import (
	"github.com/charmbracelet/log"
	"os"
	"time"

	"github.com/Inno-Gang/goodle-cli/filesystem"
	"github.com/Inno-Gang/goodle-cli/key"
	"github.com/Inno-Gang/goodle-cli/where"
)

func Init() error {
	logFile, err := filesystem.Api().OpenFile(where.LogFile(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger := log.NewWithOptions(logFile, log.Options{
		TimeFormat:      time.TimeOnly,
		ReportTimestamp: true,
		ReportCaller:    true,
	})

	level := log.ParseLevel(key.LogsLevel)
	logger.SetLevel(level)

	log.SetDefault(logger)

	return nil
}
