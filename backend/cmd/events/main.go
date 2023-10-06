package main

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spezifisch/rueder3/backend/internal/common"
	"github.com/spezifisch/rueder3/backend/pkg/events/controller"
	eventsHTTP "github.com/spezifisch/rueder3/backend/pkg/events/http"
	mockRepository "github.com/spezifisch/rueder3/backend/pkg/repository/mock"
	apiPopRepository "github.com/spezifisch/rueder3/backend/pkg/repository/pop/api"
	rabbitMQRepository "github.com/spezifisch/rueder3/backend/pkg/repository/rabbitmq"
)

func main() {
	var trustedProxies []string
	cmd := &cobra.Command{
		Use:   "events",
		Short: "HTTP Events API",
		Long:  `Rueder HTTP Events API.`,
		Run: func(cmd *cobra.Command, args []string) {
			db := common.RequireString("db")
			log.Infof("events: using pop db \"%s\"", db)

			// get options
			isDevelopmentMode := viper.GetBool("dev")
			jwtSecretKey := common.RequireString("jwt")
			if !isDevelopmentMode {
				if jwtSecretKey == "secret" || len(jwtSecretKey) < 32 {
					panic("use a JWT secret with 32 or more characters!")
				}
			}
			mqAddr := common.RequireString("rabbitmq-addr")

			// rabbitmq event source
			mqRepo := rabbitMQRepository.NewEventConsumerRepository(mqAddr)
			if mqRepo == nil {
				panic("can't connect to mq")
			}
			bind := common.RequireString("bind")
			log.Infof("events: binding to %s", bind)

			var c *controller.Controller
			if isDevelopmentMode && db == "mock" {
				c = controller.NewController(mockRepository.NewMockRepository(), mqRepo)
			} else {
				r := apiPopRepository.NewAPIPopRepository(db)
				if r == nil {
					panic("can't connect to pop repo")
				}

				c = controller.NewController(r, mqRepo)
			}

			// http server
			s := eventsHTTP.NewServer(c, bind, jwtSecretKey, isDevelopmentMode, trustedProxies)
			log.Info("🚀 events ready!")
			s.Run()

			if isDevelopmentMode {
				log.Info("❌ events quit!")
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

	cmd.PersistentFlags().StringP("bind", "b", "", "bind to ip:port")
	err = viper.BindPFlag("bind", cmd.PersistentFlags().Lookup("bind"))
	if err != nil {
		panic("BindPFlag bind failed")
	}
	viper.SetDefault("bind", ":8080")

	cmd.PersistentFlags().StringSliceVar(&trustedProxies, "trusted-proxy", []string{}, "set fiber's trusted proxy IP")
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
