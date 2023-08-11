package main

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spezifisch/rueder3/backend/internal/common"
	schedulerPopRepository "github.com/spezifisch/rueder3/backend/pkg/repository/pop/scheduler"
	"github.com/spezifisch/rueder3/backend/pkg/worker"
	"github.com/spezifisch/rueder3/backend/pkg/worker/scheduler"
)

func main() {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Feed Worker",
		Long:  `Rueder Feed Fetcher Worker.`,
		Run: func(cmd *cobra.Command, args []string) {
			db := common.RequireString("db")
			log.Infof("using pop db \"%s\"", db)

			workerCount := viper.GetInt("workers")
			if workerCount < 1 {
				workerCount = 1
			} else if workerCount > 1024 {
				workerCount = 1024
			}

			repository := schedulerPopRepository.NewSchedulerPopRepository(db)
			if repository == nil {
				return
			}

			workerPool := worker.NewFeedWorkerPool(repository)
			scheduler := scheduler.NewScheduler(repository, workerPool, workerCount)
			scheduler.Run()
		},
	}

	common.InitConfig(cmd)

	var err error
	cmd.PersistentFlags().IntP("workers", "w", 3, "worker thread count")
	err = viper.BindPFlag("workers", cmd.PersistentFlags().Lookup("workers"))
	if err != nil {
		panic(err)
	}

	cmd.PersistentFlags().Bool("dev", false, "development mode")
	err = viper.BindPFlag("dev", cmd.PersistentFlags().Lookup("dev"))
	if err != nil {
		panic(err)
	}

	err = cmd.Execute()
	if err != nil {
		log.WithError(err).Error("command failed")
	}
}
