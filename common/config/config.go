package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Mysql    MysqlConfig    `yaml:"mysql"`
	Redis    RedisConfig    `yaml:"redis"`
	RocketMQ RocketMQConfig `yaml:"rocketmq"`
	Api      API            `yaml:"api"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type MysqlConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime string `yaml:"connMaxLifetime"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type RocketMQConfig struct {
	NameSrvAddr string         `yaml:"namesrvAddr"`
	Producer    ProducerConfig `yaml:"producer"`
	Consumer    ConsumerConfig `yaml:"consumer"`
}

type ProducerConfig struct {
	Group string `yaml:"group"`
}

type ConsumerConfig struct {
	Group string `yaml:"group"`
}

type API struct {
	ThirdCall APIConfig `yaml:"thirdCall"`
	Wxd       APIConfig `yaml:"wxd"`
}

type APIConfig struct {
	EndPoint string `yaml:"endPoint"`
	Timeout  int    `yaml:"timeout"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func MustLoadConfig(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}
