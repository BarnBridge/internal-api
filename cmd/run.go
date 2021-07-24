package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/barnbridge/internal-api/api"
	"github.com/barnbridge/internal-api/config"
	"github.com/barnbridge/internal-api/db"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Say hello!",
	Long:  "Address a wonderful greeting to the majestic executioner of this CLI",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		listenOn := fmt.Sprintf(":%d", config.Store.Metrics.Port)
		sm := http.NewServeMux()
		sm.Handle("/metrics", promhttp.Handler())
		metricsSrv := &http.Server{Addr: listenOn, Handler: sm}
		go func() {
			log.Infof("serving metrics on %s", listenOn)
			err := metricsSrv.ListenAndServe()
			if err != nil && ctx.Err() == nil {
				log.Fatal(err)
			}
		}()

		db, err := db.New()
		if err != nil {
			log.Fatal(err)
		}

		a := api.New(db)
		go a.Run()

		<-ctx.Done()

		// cleanup
		_ = metricsSrv.Close()

		log.Info("Work done. Goodbye!")
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	addAPIFlags(runCmd)
	addDBFlags(runCmd)
	addAddressesFlags(runCmd)
}
