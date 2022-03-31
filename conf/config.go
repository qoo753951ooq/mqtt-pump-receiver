package conf

import (
	"log"

	"github.com/spf13/viper"
)

func Init() {

	viper.SetConfigName("config")
	viper.AddConfigPath("./conf/")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}
}
