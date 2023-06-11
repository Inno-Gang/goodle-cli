package where

import (
	"fmt"
	"github.com/Inno-Gang/goodle-cli/filesystem"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Inno-Gang/goodle-cli/app"
)

func home() string {
	home, err := os.UserHomeDir()
	if err == nil {
		return home
	}

	return "."
}

// Config path
// Will create the directory if it doesn't exist
func Config() string {
	var path string

	if customDir, present := os.LookupEnv(EnvConfigPath); present {
		return mkdir(customDir)
	}

	var userConfigDir string

	if runtime.GOOS == "darwin" {
		userConfigDir = filepath.Join(home(), ".config")
	} else {
		var err error
		userConfigDir, err = os.UserConfigDir()
		if err != nil {
			userConfigDir = filepath.Join(home(), ".config")
		}
	}

	path = filepath.Join(userConfigDir, app.Name)
	return mkdir(path)
}

// Logs path
// Will create the directory if it doesn't exist
func Logs() string {
	return mkdir(filepath.Join(Cache(), "logs"))
}

func LogFile() string {
	today := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(Logs(), fmt.Sprintf("%s.log", today))
	if !lo.Must(filesystem.Api().Exists(logFilePath)) {
		lo.Must(filesystem.Api().Create(logFilePath))
	}

	return logFilePath
}

// Cache path
// Will create the directory if it doesn't exist
func Cache() string {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		userCacheDir = "."
	}

	cacheDir := filepath.Join(userCacheDir, app.Name)
	return mkdir(cacheDir)
}

// Temp path
// Will create the directory if it doesn't exist
func Temp() string {
	tempDir := filepath.Join(os.TempDir(), app.Name)
	return mkdir(tempDir)
}
