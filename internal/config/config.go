package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"log"
	"time"
)

type Config struct {
	Env string `env:"ENVIRONMENT" envDefault:"prod"`
	HTTPServer
}

type HTTPServer struct {
	Port        string        `env:"PORT" envDefault:"8080"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" envDefault:"30s"`
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
	pflag.StringVar(&config.Port, "port", config.Port, "server's port")

	pflag.Parse()

	return config
}
