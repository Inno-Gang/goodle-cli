package main

import (
	"github.com/charmbracelet/log"
	"github.com/Inno-Gang/goodle-cli/cmd"
	"github.com/Inno-Gang/goodle-cli/config"
	"github.com/Inno-Gang/goodle-cli/logger"
	"github.com/samber/lo"
	"os"
)

func handlePanic() {
	if err := recover(); err != nil {
		log.Error("crashed", "err", err)
		os.Exit(1)
	}
}

func main() {
	defer handlePanic()

	// prepare config and logs
	lo.Must0(config.Init())
	lo.Must0(logger.Init())

	// run the app
	cmd.Execute()
}
