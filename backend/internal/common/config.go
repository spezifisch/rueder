package common

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// ConfigFile is what it is
	ConfigFile string
)

// InitConfig reads configuration from a file or environment
func InitConfig(cmd *cobra.Command) {
	var err error

	// config env/file
	cobra.OnInitialize(cobraOnInitialize)

	// config file
	cmd.PersistentFlags().StringP("config", "c", "", "config file")

	// db (-> RUEDER_DB or "db: foo" in yaml or "-d foo" or "--db foo")
	cmd.PersistentFlags().StringP("db", "d", "development|production", "pop database configuration")
	err = viper.BindPFlag("db", cmd.PersistentFlags().Lookup("db"))
	if err != nil {
		panic("BindPFlag db failed")
	}
	viper.SetDefault("db", "development")

	// log
	cmd.PersistentFlags().StringP("log", "l", "info", "log level (debug,info,warn,error,fatal)")
	err = viper.BindPFlag("log", cmd.PersistentFlags().Lookup("log"))
	if err != nil {
		panic("BindPFlag log failed")
	}
	viper.SetDefault("log", "info")

	logLevel := viper.GetString("log")
	log.SetLevel(log.MustParseLevel(logLevel))
}

func cobraOnInitialize() {
	// search path
	if ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(ConfigFile)
	} else {
		viper.AddConfigPath("/etc/rueder/")
		viper.AddConfigPath("$HOME/.rueder")
		viper.AddConfigPath(".")
	}

	// env vars
	viper.SetEnvPrefix("RUEDER")
	viper.AutomaticEnv()

	// read config
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}

// RequireString logs a fatal message if the config key doesn't exist
func RequireString(key string) (ret string) {
	ret = viper.GetString(key)
	if ret == "" {
		log.Fatalf("configuration parameter missing: %s", key)
	}
	return
}
