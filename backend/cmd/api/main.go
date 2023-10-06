package main

import (
	"os"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spezifisch/rueder3/backend/internal/common"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	ruederHTTP "github.com/spezifisch/rueder3/backend/pkg/api/http"
	mockRepository "github.com/spezifisch/rueder3/backend/pkg/repository/mock"
	apiPopRepository "github.com/spezifisch/rueder3/backend/pkg/repository/pop/api"
	rabbitMQRepository "github.com/spezifisch/rueder3/backend/pkg/repository/rabbitmq"
)

func main() {
	var trustedProxies []string
	cmd := &cobra.Command{
		Use:   "api",
		Short: "HTTP API",
		Long:  `Rueder HTTP API.`,
		Run: func(cmd *cobra.Command, args []string) {
			// first some special things for dev mode
			isDevelopmentMode := viper.GetBool("dev")

			// s6 readiness notification in dev mode
			// see: https://skarnet.org/software/s6/notifywhenup.html
			if isDevelopmentMode {
				file := os.NewFile(3, "s6ready") // hardcoded fd 3, arbitrary name
				_, err := file.Write([]byte("\n"))
				if err != nil {
					log.WithError(err).Info("s6ready not writable (this is ok on reloads)")
				} else {
					err = file.Close()
					if err != nil {
						log.WithError(err).Info("s6ready not closable")
					}
				}
			}

			// check JWT
			jwtSecretKey := common.RequireString("jwt")
			if !isDevelopmentMode {
				if jwtSecretKey == "secret" || len(jwtSecretKey) < 32 {
					panic("use a JWT secret with 32 or more characters!")
				}
			}

			// setup mq
			mqAddr := common.RequireString("rabbitmq-addr")

			// rabbitmq event sink
			mqRepo := rabbitMQRepository.NewEventPublisherRepository(mqAddr)
			if mqRepo == nil {
				panic("can't connect to mq")
			}
			go mqRepo.HandleEvents()
			defer mqRepo.Close()

			// setup SQL db
			db := common.RequireString("db")
			log.Infof("api: using pop db \"%s\"", db)

			var c *controller.Controller
			if isDevelopmentMode && db == "mock" { // allow mock sqldb only in dev mode
				c = controller.NewController(mockRepository.NewMockRepository(), mqRepo)
			} else {
				r := apiPopRepository.NewAPIPopRepository(db)
				if r == nil {
					panic("can't connect to pop repo")
				}

				c = controller.NewController(r, mqRepo)
			}

			// start http server
			s := ruederHTTP.NewServer(c, jwtSecretKey, isDevelopmentMode, trustedProxies)
			log.Info("🚀 api ready!")
			s.Run()

			if isDevelopmentMode {
				log.Info("❌ api quit!")
			}
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

	cmd.PersistentFlags().String("rabbitmq-addr", "amqp://guest:guest@rabbitmq:5672/", "RabbitMQ address")
	err = viper.BindPFlag("rabbitmq-addr", cmd.PersistentFlags().Lookup("rabbitmq-addr"))
	if err != nil {
		panic(err)
	}

	err = cmd.Execute()
	if err != nil {
		log.WithError(err).Error("command failed")
	}
}
