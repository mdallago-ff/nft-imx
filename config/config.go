package config

import (
	"log"

	"github.com/jinzhu/configor"
)

type Settings struct {
	Port                 string `default:"4000" env:"PORT"`
	AuthSecret           string `default:"" env:"AUTH_SECRET"`
	DebugMode            bool   `default:"false" env:"DEBUG"`
	AlchemyAPIKey        string `default:"" env:"ALCHEMY_API_KEY"`
	L1SignerPrivateKey   string `default:"" env:"L1_SIGNER_PRIVATE_KEY"`
	StarkPrivateKey      string `default:"" env:"STARK_PRIVATE_KEY"`
	DSN                  string `default:"" env:"DSN"`
	ProjectID            int32  `default:"0" env:"PROJECT_ID"`
	TokenDurationSeconds int64  `default:"120" env:"TOKEN_DURATION_SECONDS"`
	RedisUrl             string `default:"127.0.0.1:6379" env:"REDIS_URL"`
}

var config = Settings{}

func init() {
	if err := configor.Load(&config, "config.yml"); err != nil {
		log.Fatal(err.Error())
	}
}

func GetConfig() *Settings {
	return &config
}
