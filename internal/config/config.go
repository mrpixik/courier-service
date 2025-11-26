package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"log"
	"os"
	"time"
)

type Config struct {
	Env string //local, dev, prod ||| default: local
	//Env      string          `env:"ENVIRONMENT,required"` //local, dev, prod -- пока применяется только для настройки логгера
	Postgres                   PostgresStorage `envPrefix:"POSTGRES_"`
	HTTP                       HTTPServer      `envPrefix:"HTTP_"`
	DeliveryWorkerTickInterval time.Duration   `env:"DELIVERY_WORKER_TICK_INTERVAL" envDefault:"60s"`
}

type HTTPServer struct {
	Port            string        `env:"PORT" envDefault:"8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"30s"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
}

type PostgresStorage struct {
	User            string        `env:"USER,required"`
	Password        string        `env:"PASSWORD,required"`
	Host            string        `env:"HOST,required"`
	Port            string        `env:"PORT,required"`
	DbName          string        `env:"DB_NAME,required"`
	MaxConns        int32         `env:"MAX_CONNS,required"`
	MinConns        int32         `env:"MIN_CONNS,required"`
	MaxConnLifeTime time.Duration `env:"MAX_CONN_LIFE_TIME,required"`
}

func MustLoad() Config {

	environment := os.Getenv("ENVIRONMENT")
	var httpPort string
	if environment == "" {
		pflag.StringVar(&environment, "env", environment, "environment (local, prod)")
		pflag.StringVar(&httpPort, "port", httpPort, "server's port")
		pflag.Parse()
	}
	fmt.Printf(".%s.\n", environment)
	if environment != "prod" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("unable to load config from .env: \n%s", err.Error())
		}
	}
	var config Config

	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("unable to load config: \n%s", err.Error())
	}

	config.Env = environment

	// Значение порта переопределяется только в случае, если в --port передается какое-то значение
	if httpPort != "" {
		config.HTTP.Port = httpPort
	}

	return config
}
