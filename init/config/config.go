package config

import "github.com/spf13/viper"

var ServerConfig Config

type Config struct {
	Debug        bool     `mapstructure:"debug"`
	Port         int      `mapstructure:"port"`
	Entry        string   `mapstructure:"entry"`
	PasswordSalt string   `mapstructure:"password_salt"`
	Postgres     postgres `mapstructure:"db"`
	Log          log      `mapstructure:"log"`
}

type postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type log struct {
	HttpLoggerPath     string `mapstructure:"http_logger_path"`
	PostgresLoggerPath string `mapstructure:"postgres_logger_path"`
	QuizLoggerPath     string `mapstructure:"quiz_logger_path"`
}

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&ServerConfig); err != nil {
		return err
	}

	return nil
}
