package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

//config contains all global config of application
type config struct {
	Env     string
	Common  CommonConfig
	DB      DBConfig
	Http    HttpConfig
	Telebot TelebotConfig
	Binance BinanceConfig
}

var cfg *config

//Instance return instance of global config
func Instance() *config {
	if cfg == nil {
		cfg = &config{
			Common:  CommonConfig{},
			DB:      DBConfig{},
			Http:    HttpConfig{},
			Telebot: TelebotConfig{},
			Binance: BinanceConfig{},
		}
	}
	return cfg
}

//Load loads configurations from file and env
func Load(configFile string) error {
	// Default config values
	c := Instance()
	defaults.SetDefaults(c)

	// --- hacking to load reflect structure config into env ----//
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Read config file failed. ", err)

		configBuffer, err := json.Marshal(c)

		if err != nil {
			return err
		}

		err = viper.ReadConfig(bytes.NewBuffer(configBuffer))
		if err != nil {
			return err
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// -- end of hacking --//

	fmt.Println(viper.GetString("ENV"))
	viper.AutomaticEnv()
	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
