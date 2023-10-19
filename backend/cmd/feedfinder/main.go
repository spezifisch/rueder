package main

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spezifisch/rueder3/backend/internal/common"
	"github.com/spezifisch/rueder3/backend/pkg/feedfinder/controller"
	feedfinderHTTP "github.com/spezifisch/rueder3/backend/pkg/feedfinder/http"
)

func main() {
	var trustedProxies []string
	cmd := &cobra.Command{
		Use:   "feedfinder",
		Short: "Feedfinder Service",
		Long:  `Rueder Feedfinder parses websites and returns their feeds.`,
		Run: func(cmd *cobra.Command, args []string) {
			isDevelopmentMode := viper.GetBool("dev")
			jwtSecretKey := common.RequireString("jwt")
			if !isDevelopmentMode {
				if jwtSecretKey == "secret" || len(jwtSecretKey) < 32 {
					panic("use a JWT secret with 32 or more characters!")
				}
			}
			bind := common.RequireString("bind")
			log.Infof("feedfinder: binding to %s", bind)

			c := controller.NewController()
			s := feedfinderHTTP.NewServer(c, bind, jwtSecretKey, isDevelopmentMode, trustedProxies)
			log.Info("🚀 feedfinder ready!")
			s.Run()

			if isDevelopmentMode {
				log.Info("❌ feedfinder quit!")
			}
		},
	}

	common.InitConfig(cmd)

	var err error
	cmd.PersistentFlags().Bool("dev", false, "development mode")
	err = viper.BindPFlag("dev", cmd.PersistentFlags().Lookup("dev"))
	if err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringP("bind", "b", "", "bind to ip:port")
	err = viper.BindPFlag("bind", cmd.PersistentFlags().Lookup("bind"))
	if err != nil {
		panic("BindPFlag bind failed")
	}
	viper.SetDefault("bind", ":8080")

	// log
	cmd.PersistentFlags().StringSliceVar(&trustedProxies, "trusted-proxy", []string{}, "set fiber's trusted proxy IP")
	err = viper.BindPFlag("trusted-proxy", cmd.PersistentFlags().Lookup("trusted-proxy"))
	if err != nil {
		panic(err)
	}

	err = cmd.Execute()
	if err != nil {
		log.WithError(err).Error("command failed")
	}
}
