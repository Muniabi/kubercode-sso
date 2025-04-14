package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Port    int `yaml:"port"`
	Timeout int `yaml:"timeout"`
}

type Config struct {
	Env                          string     `yaml:"env" env-default:"local"`
	HTTP                         HTTPConfig `yaml:"http"`
	MigrationPath                string     `yaml:"migrationPath" env-default:"migrations"`
	PrivateKeyPath               string     `yaml:"privateKeyPath" env-default:"certs/jwtRSA256-private.pem"`
	PublicKeyPath                string     `yaml:"publicKeyPath" env-default:"certs/jwtRSA256-public.pem"`
	JWTAlgorithm                 string     `yaml:"JWTAlgorithm" env-default:""`
	RefreshTokenDurationDays     int        `yaml:"refreshTokenDurationDays"`
	AccessTokenDurationMinutes   int        `yaml:"accessTokenDurationMinutes"`
	EventStoreConnectionString   string     `yaml:"eventStoreConnectionString"`
	MongoDBConnectionString      string     `yaml:"mongoDbConnectionString"`
	ProjectionsGroupName         string     `yaml:"projectionsGroupName"`
	AccountPrefix                string     `yaml:"accountPrefix"`
	RedisAddress                 string     `yaml:"redisAddress"`
	RedisPassword                string     `yaml:"redisPassword"`
	OTPExpirationDurationSeconds int        `yaml:"otpExpirationDurationSeconds"`
	OTPLength                    int        `yaml:"otpLength"`
	SMTPHost                     string     `yaml:"smtpHost"`
	SMTPPort                     string     `yaml:"smtpPort"`
	SMPTUsername                 string     `yaml:"smtpUsername"`
	SMTPPassword                 string     `yaml:"SMPTPassword"`
	OurEmail                     string     `yaml:"ourEmail"`
	SMTPGmailHost                string     `yaml:"smtpGmailHost"`
	SMTPGmailPort                string     `yaml:"smtpGmailPort"`
	GmailEmail                   string     `yaml:"gmailEmail"`
	NatsURL                      string     `yaml:"natsURL"`
	Subject                      string     `yaml:"subject"`
}

func fetchConfigPath(filename string) string {
	path, err := filepath.Abs(filename)
	if err != nil {
		errString := fmt.Sprintf("Error getting absolute path of config file %s: %s\n", filename, err)
		panic(errString)
	}
	return path
}

func MustLoadConfig(filename string) *Config {
	configPath := fetchConfigPath(filename)
	if configPath == "" {
		panic("config path not set")
	}
	fmt.Println(configPath)
	if q, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println(q)
		panic("config path does not exist")
	}
	config := &Config{}
	if err := cleanenv.ReadConfig(configPath, config); err != nil {
		panic(err)
	}
	return config
}
