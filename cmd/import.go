package cmd

import (
	"github.com/koolay/slow-api/app"
	"github.com/koolay/slow-api/config"
	"github.com/koolay/slow-api/logging"
	"github.com/spf13/cobra"
)

var (
	file string
)

// serveCmd represents the serve command
var inputCmd = &cobra.Command{
	Use:   "import",
	Short: "import slow log to storage",
	Run: func(cmd *cobra.Command, args []string) {
		collector := app.NewMySqlCollector(config.Context)
		logging.Logger.INFO.Println("Import mysql slow log: " + file)
		if err := collector.ImportLogFile(file); err != nil {
			logging.Logger.ERROR.Println(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(inputCmd)
	inputCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Path of slow log")
}
