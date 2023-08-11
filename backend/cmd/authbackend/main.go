package main

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spezifisch/rueder3/backend/internal/common"
	"github.com/spezifisch/rueder3/backend/pkg/authbackend/controller"
	authBackendHTTP "github.com/spezifisch/rueder3/backend/pkg/authbackend/http"
	authBackendPopRepository "github.com/spezifisch/rueder3/backend/pkg/repository/pop/authbackend"
)

func main() {
	cmd := &cobra.Command{
		Use:   "authbackend",
		Short: "Authentication Backend",
		Long:  `Rueder Authentication Backend provides auth claims for loginsrv.`,
		Run: func(cmd *cobra.Command, args []string) {
			db := common.RequireString("db")
			log.Infof("using pop db \"%s\"", db)

			isDevelopmentMode := viper.GetBool("dev")

			r := authBackendPopRepository.NewAuthBackendPopRepository(db)
			if r == nil {
				return
			}
			bind := common.RequireString("bind")
			log.Infof("authbackend: binding to %s", bind)

			c := controller.NewController(r)
			s := authBackendHTTP.NewServer(c, bind, isDevelopmentMode)
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

	cmd.PersistentFlags().StringP("bind", "b", "", "bind to ip:port")
	err = viper.BindPFlag("bind", cmd.PersistentFlags().Lookup("bind"))
	if err != nil {
		panic("BindPFlag bind failed")
	}
	viper.SetDefault("bind", ":8080")

	err = cmd.Execute()
	if err != nil {
		log.WithError(err).Error("command failed")
	}
}
