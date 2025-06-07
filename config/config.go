package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogLevel zapcore.Level `env:"LOG_LEVEL" envDefault:"debug"`
	JsonLogs bool          `env:"JSON_LOGS" envDefault:"false"`
	OneShot  bool          `env:"ONE_SHOT" envDefault:"false"`

	ClientId string `env:"BOT_ID" envDefault:"668946506836869150"`

	Auth struct {
		TopGG          string `env:"TOPGG" envDefault:""`
		DiscordBotList string `env:"DISCORDBOTLIST" envDefault:""`
		BotListMe      string `env:"BOTLISTME" envDefault:""`
	} `envPrefix:"AUTH_"`

	DatabaseUrl string `env:"DATABASE_URL,required"`
}

var Conf Config

func Parse() {
	var err error
	if _, err = os.Stat(".env"); err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}

	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
