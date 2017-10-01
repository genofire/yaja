package cmd

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Runs the yaja server",
	Example: "yaja serve -config /etc/yaja.toml",
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		// Wait for INT/TERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Infoln("received", sig)

	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "config.toml", "Path to configuration file")
}
