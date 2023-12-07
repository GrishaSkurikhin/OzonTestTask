package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Server struct {
	RestPort     int
	GRPCPort     int
	ShortURLHost string
	Env          string
	DB           Database
	InMemory     bool
}

type Database struct {
	Source string
}

func MustLoad() Server {
	viper.SetDefault("env", "local")
	viper.SetDefault("short_url_host", "http://localhost:8080")

	viper.BindEnv("env")
	viper.BindEnv("short_url_host")
	viper.BindEnv("db_source")

	viper.MustBindEnv("rest_port")
	viper.MustBindEnv("grpc_port")

	pflag.Bool("memory", false, "use in-memory source")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	return Server{
		RestPort:     viper.GetInt("rest_port"),
		GRPCPort:     viper.GetInt("grpc_port"),
		ShortURLHost: viper.GetString("short_url_host"),
		Env:          viper.GetString("env"),
		InMemory:     viper.GetBool("memory"),
		DB: Database{
			Source: viper.GetString("db_source"),
		},
	}
}
