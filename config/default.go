package config

import (
	"github.com/Inno-Gang/goodle-cli/key"
	"github.com/spf13/viper"
)

// fields is the config fields with their default values and descriptions
var fields = []*Field{
	// LOGS
	{
		key.LogsWrite,
		true,
		"Write logs to file",
	},
	{
		key.LogsLevel,
		"info",
		`Logs level.
Available options are: (from less to most verbose)
fatal, error, warn, info, debug`,
	},
	// END LOGS

	// AUTH
	{
		key.AuthEmail,
		"",
		"Innopolis email",
	},
	{
		key.AuthPassword,
		"",
		"Innopolis password in a plaintext. Not recommended to use",
	},
	{
		key.AuthRemember,
		false,
		"Remember password. Warning - passwords are stored in a plaintext",
	},
	// END AUTH

	// TUI
	{
		key.TUIShowSections,
		false,
		"Show sections as a separate state",
	},
	// END TUI
}

func setDefaults() {
	Default = make(map[string]*Field, len(fields))
	for _, f := range fields {
		Default[f.Key] = f
		viper.SetDefault(f.Key, f.DefaultValue)
		viper.MustBindEnv(f.Key)
	}
}

var Default map[string]*Field
