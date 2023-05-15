package logger

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"

	"github.com/Inno-Gang/goodle-cli/filesystem"
	"github.com/Inno-Gang/goodle-cli/key"
	"github.com/Inno-Gang/goodle-cli/where"
	"github.com/samber/lo"
)

func Init() error {
	logsPath := where.Logs()

	if logsPath == "" {
		return errors.New("logs path is not set")
	}

	today := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logsPath, fmt.Sprintf("%s.log", today))
	if !lo.Must(filesystem.Api().Exists(logFilePath)) {
		lo.Must(filesystem.Api().Create(logFilePath))
	}

	logFile, err := filesystem.Api().OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger := log.NewWithOptions(logFile, log.Options{
		TimeFormat:      time.TimeOnly,
		ReportTimestamp: true,
		ReportCaller:    viper.GetBool(key.LogsReportCaller),
	})

	level := log.ParseLevel(key.LogsLevel)
	logger.SetLevel(level)

	log.SetDefault(logger)

	return nil
}
