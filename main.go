package main

import (
	"github.com/Inno-Gang/goodle-cli/cache"
	"github.com/Inno-Gang/goodle-cli/cmd"
	"github.com/Inno-Gang/goodle-cli/config"
	"github.com/Inno-Gang/goodle-cli/logger"
	"github.com/samber/lo"
)

func main() {
	defer cache.Close()

	// prepare config and logs
	lo.Must0(config.Init())
	lo.Must0(logger.Init())

	// run the app
	cmd.Execute()
}
