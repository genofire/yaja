package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"dev.sum7.eu/genofire/yaja/model/config"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var configPath string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Runs the yaja server",
	Example: "yaja serve -c /etc/yaja.conf",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := config.ReadConfigFile(configPath)
		if err != nil {
			log.Fatal("unable to load config file:", err)
		}

		log.Infoln("yaja started ")

		// Wait for INT/TERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Infoln("received", sig)

	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "yaja.conf", "Path to configuration file")
}
