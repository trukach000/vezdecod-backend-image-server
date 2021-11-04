package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

// Default config for local environment
type Config struct {
	Host   string `env:"HOST" envDefault:"127.0.0.1"`
	Port   string `env:"PORT" envDefault:"80"`
	UseSSL bool   `env:"USE_SSL" envDefault:"false"`

	DBHost string `env:"DB_HOST" envDefault:"localhost"`
	DBPort string `env:"DB_PORT" envDefault:"3306"`
	DBUser string `env:"DB_USER" envDefault:"imloader"`
	DBPass string `env:"DB_PASS" envDefault:"password"`
	DBName string `env:"DB_NAME" envDefault:"imloader"`

	TarantoolHost string `env:"TARANTOOL_HOST" envDefault:"localhost"`
	TarantoolPort string `env:"TARANTOOL_PORT" envDefault:"3301"`
	TarantoolUser string `env:"TARANTOOL_USER" envDefault:"admin"`
	TarantoolPass string `env:"TARANTOOL_PASS" envDefault:"pass"`

	GlobalPrefix string `env:"GLOBAL_PREFIX" envDefault:"imloader"`

	SecretAuthKey string `env:"SECRET_KEY" envDefault:"123"`
}

var conf = Config{}

func Init() {
	err := env.Parse(&conf)
	if err != nil {
		logrus.Error("failed to parse config from env: %s", err)
		panic(err)
	}
}

func Get() Config {
	return conf
}
