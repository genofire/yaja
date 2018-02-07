package cmd

import (
	"dev.sum7.eu/genofire/yaja/daemon"
	"github.com/spf13/cobra"
)

// DaemonCMD represents the serve command
var DaemonCMD = &cobra.Command{
	Use:   "daemon",
	Short: "daemon of yaja",
}

func init() {
	DaemonCMD.AddCommand(daemon.ServerCMD)
	DaemonCMD.AddCommand(daemon.TesterCMD)
	RootCMD.AddCommand(DaemonCMD)
}
