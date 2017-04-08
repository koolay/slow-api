package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/koolay/slow-api/app"
	"github.com/koolay/slow-api/config"
	"github.com/koolay/slow-api/logging"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as serve",
	Run: func(cmd *cobra.Command, args []string) {
		signalChan := make(chan os.Signal, 1)
		doneChan := make(chan bool)
		errChan := make(chan error, 10)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			collector := app.NewMySqlCollector(config.Context)
			logging.Logger.INFO.Println("start collect mysql slow log")
			if err := collector.Start(); err != nil {
				logging.Logger.ERROR.Println(err.Error())
			}
		}()

		for {
			select {
			case err := <-errChan:
				logging.Logger.ERROR.Printf("%v", err)
			case s := <-signalChan:
				logging.Logger.INFO.Printf("Captured %v. Exiting...", s)
				close(doneChan)
			case <-doneChan:
				logging.Logger.DEBUG.Println("Done!")
				os.Exit(0)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

}
