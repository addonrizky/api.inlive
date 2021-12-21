package pkg

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`.env`)
	viper.SetConfigType(`env`)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
	}
}

func GetConfigString(key string) string {
	return viper.GetString(fmt.Sprintf("%v", key))
}

func GetConfigInt(key string) int {
	return viper.GetInt(fmt.Sprintf("%v", key))
}

func GetConfigDuration(key string) time.Duration {
	return viper.GetDuration(fmt.Sprintf("%v", key))
}

func GetConfigBool(key string) bool {
	return viper.GetBool(fmt.Sprintf("%v", key))
}
