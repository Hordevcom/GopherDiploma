package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAdress          string `env:"SERVER_ADDRESS"`
	DatabaseDsn           string `env:"DATABASE_DSN"`
	AccurualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() Config {
	var conf Config
	err := env.Parse(&conf)

	if err != nil {
		fmt.Println(err)
	}

	if conf.DatabaseDsn != "" && conf.ServerAdress != "" && conf.AccurualSystemAddress != "" {
		return conf
	}

	if conf.DatabaseDsn == "" {
		flag.StringVar(&conf.DatabaseDsn, "d", "", "database dsn") //"postgres://postgres:1@localhost:5432/postgres"
	}

	if conf.AccurualSystemAddress == "" {
		flag.StringVar(&conf.AccurualSystemAddress, "r", "", "ACCRUAL_SYSTEM_ADDRESS")
	}

	flag.StringVar(&conf.ServerAdress, "a", "localhost:8080", "server address")

	return conf
}
