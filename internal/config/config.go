package config

import (
	"log"
	"sync"
	"time"

	auth "auth-api/pkg/auth/jwt"
	"auth-api/pkg/client/mysql"
	"auth-api/pkg/email/smtp"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                     string `yaml:"env" env:"ENV" env-default:"local"`
	mysql.StorageConfig     `yaml:"storage"`
	HTTPServerConfig        `yaml:"http_server"`
	smtp.SenderConfig       `yaml:"email_sender"`
	auth.TokenManagerConfig `yaml:"jwt_token"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
	User        string        `yaml:"user" env:"HTTP_SERVER_USER" env-required:"true"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Println("Read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("..\\config\\config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instance
}
