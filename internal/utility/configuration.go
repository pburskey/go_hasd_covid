package utility

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Configuration struct {
	Env   string
	Redis RedisConfiguration
	MySQL map[string]interface{}
}

type RedisConfiguration struct {
	Password string
}

func LoadConfiguration(path string) *Configuration {
	viper.SetConfigName("development") // config file name without extension
	viper.SetConfigType("yaml")
	//viper.AddConfigPath(".")

	//viper.AddConfigPath("./config/") // config file path
	viper.AddConfigPath(filepath.Join(path, "/config/"))
	viper.AutomaticEnv() // read value ENV variable

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	// Set default value
	viper.SetDefault("app.linetoken", "DefaultLineTokenValue")

	var config *Configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	// Declare var
	env := viper.GetString("env")

	// Print
	fmt.Println("---------- Example ----------")
	fmt.Println("app.env :", env)

	return config

}
