package handlers

import (
	"github.com/spf13/viper"
)

var DBconfig *DatabasesConfig

type DatabasesConfig struct {
	Database []struct {
		Driver   string `json:"driver"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		DbName   string `json:"dbName"`
		Schema   string `json:"schema,omitempty"`
		Tables   []struct {
			Select string `json:"select"`
			TbName string `json:"tbName"`
			Where  string `json:"where,omitempty"`
		} `json:"tables,omitempty"`
	} `json:"database"`
}

func LoadConfig() {
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name

	Logger.Sugar().Infof("%s", "-----------------------------------------------")
	Logger.Sugar().Infof("%s", "JSON Config file")
	Logger.Sugar().Infof("%s", "-----------------------------------------------")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			Logger.Sugar().Fatalf("[%s] %s - %s", "X", "config.json", err)
		} else {
			// Config file was found but another error was produced
			Logger.Sugar().Fatalf("[%s] %s - %s", "X", "config.json", err)
		}
	}

	if err := viper.Unmarshal(&DBconfig); err != nil {
		Logger.Sugar().Fatalf("[%s] %s - %s", "X", "config.json", err)
	} else {
		Logger.Sugar().Infof("[%s] %s", "/", "config.json")
	}
}
