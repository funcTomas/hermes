package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig         `yaml:"server"`
	Mysql    MysqlConfig          `yaml:"mysql"`
	Redis    RedisConfig          `yaml:"redis"`
	RocketMQ RocketMQConfig       `yaml:"rocketmq"`
	API      map[string]APIConfig `yaml:"api"`
}

type APIConfig struct {
	EndPoint string `yaml:"endPoint"`
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
	Topic string `yaml:"topic"`
}

type ConsumerConfig struct {
	Group  string        `yaml:"group"`
	Topics []TopicConfig `yaml:"topics"`
}

type TopicConfig struct {
	Name string `yaml:"name"`
	// Selector MessageSelector `yaml:"selector"` // 如果需要，可以加上
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
