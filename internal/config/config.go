package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Server struct {
	Address  string
	Env      string
	DB       Database
	InMemory bool
}

type Database struct {
	Type   string
	Source string
}

func MustLoad() Server {
	viper.BindEnv("loglevel")
	viper.SetDefault("env", "local")

	viper.MustBindEnv("server_address")
	viper.MustBindEnv("db_type")
	viper.MustBindEnv("db_source")

	pflag.Bool("memory", false, "use in-memory source")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	return Server{
		Address:  viper.GetString("server_address"),
		Env:      viper.GetString("env"),
		InMemory: viper.GetBool("memory"),
		DB: Database{
			Type:   viper.GetString("db_type"),
			Source: viper.GetString("db_source"),
		},
	}
}
