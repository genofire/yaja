package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	timestamps bool
)

// RootCmd represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "yaja",
	Short: "Yet another jabber implementation",
	Long:  `A small standalone command line round about jabber (e.g tester WIP: client & server)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCMD.Execute(); err != nil {
		log.Panicln(err)
	}
}
