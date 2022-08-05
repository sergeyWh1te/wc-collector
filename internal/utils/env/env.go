package env

import (
	"context"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	AppConfig AppConfig
	PgConfig  PgConfig
}

type AppConfig struct {
	InfuraKey              string
	DepositContractAddress string
	WC                     string
	InfuraProjectID        string
	InfuraProjectSecret    string
	InfuraETH2Host         string
	FromBlock              int64
	ToBlock                int64
	Step                   int64
}

type PgConfig struct {
	Port     uint
	Host     string
	Username string
	Password string
	Database string
	Schema   string
	SslMode  string
}

var (
	cfg Config

	onceDefaultClient sync.Once
)

func Read(ctx context.Context) (*Config, error) {
	var err error

	onceDefaultClient.Do(func() {
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.SetConfigFile(".env")

		viper.AutomaticEnv()
		if viperErr := viper.ReadInConfig(); err != nil {
			if _, ok := viperErr.(viper.ConfigFileNotFoundError); !ok {
				err = viperErr
				return
			}
		}

		cfg = Config{
			AppConfig: AppConfig{
				InfuraKey:              viper.GetString("INFURA_KEY"),
				DepositContractAddress: viper.GetString("DEPOSIT_CONTRACT_ADDRESS"),
				WC:                     viper.GetString("WITHDRAWAL_CREDENTIALS"),
				InfuraProjectID:        viper.GetString("INFURA_PROJECT_ID"),
				InfuraProjectSecret:    viper.GetString("INFURA_PROJECT_SECRET"),
				InfuraETH2Host:         viper.GetString("INFURA_EHT2_HOST"),

				FromBlock: viper.GetInt64("FROM_BLOCK"),
				ToBlock:   viper.GetInt64("TO_BLOCK"),
				Step:      viper.GetInt64("STEP"),
			},
			PgConfig: PgConfig{
				Port:     viper.GetUint("PG_PORT"),
				Host:     viper.GetString("PG_HOST"),
				Username: viper.GetString("PG_USERNAME"),
				Password: viper.GetString("PG_PASSWORD"),
				Database: viper.GetString("PG_DATABASE"),
				Schema:   viper.GetString("PG_SCHEMA"),
				SslMode:  viper.GetString("PG_SSLMODE"),
			},
		}
	})

	return &cfg, err
}
