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

			c := controller.NewController()
			s := feedfinderHTTP.NewServer(c, jwtSecretKey, isDevelopmentMode, trustedProxies)
			s.Run()
		},
	}

	common.InitConfig(cmd)

	var err error
	cmd.PersistentFlags().Bool("dev", false, "development mode")
	err = viper.BindPFlag("dev", cmd.PersistentFlags().Lookup("dev"))
	if err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringSliceVar(&trustedProxies, "trusted-proxy", []string{}, "set gin's trusted proxy IP")
	err = viper.BindPFlag("trusted-proxy", cmd.PersistentFlags().Lookup("trusted-proxy"))
	if err != nil {
		panic(err)
	}

	err = cmd.Execute()
	if err != nil {
		log.WithError(err).Error("command failed")
	}
}
