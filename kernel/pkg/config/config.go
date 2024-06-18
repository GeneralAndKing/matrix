package config

import (
	"github.com/spf13/viper"
)

func UnmarshalConfig(config interface{}, appName string, configName string) error {
	viper.SetEnvPrefix(appName)
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")

	viper.SetConfigName(configName)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.SetConfigName(extendConfigName(configName))
	_ = viper.MergeInConfig()
	if err := viper.Unmarshal(config); err != nil {
		return err
	}
	return nil
}

func extendConfigName(configName string) string {
	if viper.GetString("env") == "" {
		return configName
	} else {
		return configName + "-" + viper.GetString("env")
	}
}

func init() {
	viper.AutomaticEnv()
}
