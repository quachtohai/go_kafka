package config

import (
	"flag"
	"fmt"
	"go_kafka/pkg/constants"
	kafkaClient "go_kafka/pkg/kafka"
	"go_kafka/pkg/logger"
	"os"

	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Reader microservice config path")
}

type EnvConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBName     string `env:"DB_NAME,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`
}
type KafkaTopics struct {
	EventCreated kafkaClient.TopicConfig `mapstructure:"eventCreated"`
	EventUpdated kafkaClient.TopicConfig `mapstructure:"eventUpdated"`
	EventDeleted kafkaClient.TopicConfig `mapstructure:"eventDeleted"`
}
type KafkaConfig struct {
	Logger      *logger.Config      `mapstructure:"logger"`
	ServiceName string              `mapstructure:"serviceName"`
	KafkaTopics KafkaTopics         `mapstructure:"kafkaTopics"`
	Kafka       *kafkaClient.Config `mapstructure:"kafka"`
}

func NewEnvConfig() *EnvConfig {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Unable to load .env: %e", err)
	}

	config := &EnvConfig{}

	if err := env.Parse(config); err != nil {
		log.Fatalf("Unable to load variables from .env: %e", err)
	}

	return config
}
func InitConfig() (*KafkaConfig, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")

			}
			configPath = fmt.Sprintf("%s/config/config.yaml", getwd)
		}
	}

	cfg := &KafkaConfig{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}

	return cfg, nil
}
