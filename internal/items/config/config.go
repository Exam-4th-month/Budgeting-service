package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server        ServerConfig
		MongoDb      MongoDbConfig
		Redis         RedisConfig
		JWT           JWTConfig
		RabbitMQ      RabbitMQConfig
	}
	JWTConfig struct {
		SecretKey string
	}

	ServerConfig struct {
		Port string
	}
	MongoDbConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	RedisConfig struct {
		Host string
		Port string
	}
	RabbitMQConfig struct {
		RabbitMQ string
	}
)

func (c *Config) Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	c.Server.Port = ":" + os.Getenv("SERVER_PORT")
	c.MongoDb.Host = os.Getenv("DB_HOST")
	c.MongoDb.Port = os.Getenv("DB_PORT")
	c.MongoDb.User = os.Getenv("DB_USER")
	c.MongoDb.Password = os.Getenv("DB_PASSWORD")
	c.MongoDb.DBName = os.Getenv("DB_NAME")
	c.Redis.Host = os.Getenv("REDIS_HOST")
	c.Redis.Port = os.Getenv("REDIS_PORT")
	c.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	c.RabbitMQ.RabbitMQ = os.Getenv("RABBITMQ_URI")

	return nil
}

func New() (*Config, error) {
	var config Config
	if err := config.Load(); err != nil {
		return nil, err
	}
	return &config, nil
}
