package main

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spezifisch/rueder3/backend/internal/common"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	ruederHTTP "github.com/spezifisch/rueder3/backend/pkg/api/http"
	mockRepository "github.com/spezifisch/rueder3/backend/pkg/repository/mock"
	apiPopRepository "github.com/spezifisch/rueder3/backend/pkg/repository/pop/api"
)

func main() {
	var trustedProxies []string
	cmd := &cobra.Command{
		Use:   "api",
		Short: "HTTP API",
		Long:  `Rueder HTTP API.`,
		Run: func(cmd *cobra.Command, args []string) {
			db := common.RequireString("db")
			log.Infof("using pop db \"%s\"", db)

			isDevelopmentMode := viper.GetBool("dev")
			jwtSecretKey := common.RequireString("jwt")
			if !isDevelopmentMode {
				if jwtSecretKey == "secret" || len(jwtSecretKey) < 32 {
					panic("use a JWT secret with 32 or more characters!")
				}
			}

			var c *controller.Controller
			if db != "mock" {
				r := apiPopRepository.NewAPIPopRepository(db)
				if r == nil {
					return
				}

				c = controller.NewController(r)
			} else {
				c = controller.NewController(mockRepository.NewMockRepository())
			}

			s := ruederHTTP.NewServer(c, jwtSecretKey, isDevelopmentMode, trustedProxies)
			s.Run()
		},
	}

	common.InitConfig(cmd)

	var err error
	cmd.PersistentFlags().String("jwt", "", "JWT secret key")
	err = viper.BindPFlag("jwt", cmd.PersistentFlags().Lookup("jwt"))
	if err != nil {
		panic(err)
	}

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
