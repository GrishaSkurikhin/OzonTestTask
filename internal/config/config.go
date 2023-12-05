package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type server struct {
	Address  string
	LogLevel string
	DB       database
	InMemory bool
}

type database struct {
	Type   string
	Source string
}

func MustLoad() server {
	viper.BindEnv("loglevel")
	viper.SetDefault("loglevel", "debug")

	viper.MustBindEnv("server_address")
	viper.MustBindEnv("db_type")
	viper.MustBindEnv("db_source")

	pflag.Bool("memory", false, "use in-memory source")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	return server{
		Address:  viper.GetString("server_address"),
		LogLevel: viper.GetString("loglevel"),
		InMemory: viper.GetBool("memory"),
		DB: database{
			Type:   viper.GetString("db_type"),
			Source: viper.GetString("db_source"),
		},
	}
}
