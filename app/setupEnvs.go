package app

import "github.com/spf13/viper"

func GetEnvs(key string) string {
	viper.SetConfigFile(".env")
    viper.ReadInConfig()

    env := viper.Get(key).(string)

	return env
	
}