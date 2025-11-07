package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"log"
	"time"
)

type Config struct {
	Env string `env:"ENVIRONMENT" envDefault:"prod"` //local, dev, prod -- пока применяется только для настройки логгера
	PostgresStorage
	HTTPServer
}

type HTTPServer struct {
	ServerPort  string        `env:"PORT" envDefault:"8080"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" envDefault:"30s"`
}

type PostgresStorage struct {
	User            string        `env:"POSTGRES_USER,required"`
	Password        string        `env:"POSTGRES_PASSWORD,required"`
	Host            string        `env:"POSTGRES_HOST,required"`
	DbPort          string        `env:"POSTGRES_PORT,required"`
	DbName          string        `env:"POSTGRES_DB_NAME,required"`
	MaxConns        int32         `env:"POSTGRES_MAX_CONNS,required"`
	MinConns        int32         `env:"POSTGRES_MIN_CONNS,required"`
	MaxConnLifeTime time.Duration `env:"POSTGRES_MAX_CONN_LIFE_TIME,required"`
}

func MustLoad() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load config: \n%s", err.Error())
	}

	var config Config

	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("unable to load config: \n%s", err.Error())
	}

	// Значение порта переопределяется только в случае, если в --port передается какое-то значение
	pflag.StringVar(&config.ServerPort, "port", config.ServerPort, "server's port")

	pflag.Parse()

	return config
}
