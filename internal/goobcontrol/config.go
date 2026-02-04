package goobcontrol

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

func CreateConfig() *viper.Viper {
	botConfig := viper.New()
	botConfig.SetEnvPrefix("goob")
	botConfig.SetTypeByDefaultValue(true)
	botConfig.SetConfigName("goobconfig")
	botConfig.AddConfigPath(".")
	if err := botConfig.ReadInConfig(); err != nil {
		var configFileNotfoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotfoundError) {
			//Ignore these errors as we might be reading from env
		} else {
			panic("Couldn't load config file: " + err.Error())
		}
	}
	botConfig.SetDefault("bot.name", "Goob Control")
	botConfig.SetDefault("bot.debug", false)
	botConfig.SetDefault("discord.privateGuilds", []string{})
	botConfig.SetDefault("database.insecure", false)
	botConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	botConfig.AutomaticEnv()
	return botConfig
}
