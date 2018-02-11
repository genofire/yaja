package daemon

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/golang-lib/file"
	"dev.sum7.eu/genofire/golang-lib/worker"
	"dev.sum7.eu/genofire/yaja/client"
	"dev.sum7.eu/genofire/yaja/daemon/tester"
	"dev.sum7.eu/genofire/yaja/messages"

	"github.com/spf13/cobra"
)

var (
	configTester   = &tester.Config{}
	testerInstance = tester.NewTester()
	testerWorker   *worker.Worker
)

// TesterCMD represents the serve command
var TesterCMD = &cobra.Command{
	Use:     "tester",
	Short:   "runs xmpp tester server",
	Example: "yaja daemon tester -c /etc/yaja.conf",
	Run: func(cmd *cobra.Command, args []string) {

		if err := file.ReadTOML(configPath, configTester); err != nil {
			log.Fatal("unable to load config file:", err)
		}

		log.SetLevel(configTester.Logging)

		if err := file.ReadJSON(configTester.AccountsPath, testerInstance); err != nil {
			log.Warn("unable to load state file:", err)
		}

		mainClient, err := client.NewClientProtocolDuration(configTester.Client.JID, configTester.Client.Password, "tcp", configTester.Timeout.Duration)
		if err != nil {
			log.Fatal("unable to connect with main jabber client: ", err)
		}
		defer mainClient.Close()

		for _, admin := range configTester.Admins {
			mainClient.Send(&messages.MessageClient{
				To:   admin,
				Type: "chat",
				Body: "yaja tester starts",
			})
		}
		testerInstance.Timeout = configTester.Timeout.Duration
		testerInstance.Start(mainClient, configTester.Client.Password)
		testerInstance.CheckStatus()
		testerWorker = worker.NewWorker(configTester.Interval.Duration, func() {
			testerInstance.CheckStatus()
			file.SaveJSON(configTester.AccountsPath, testerInstance)
			file.SaveJSON(configTester.OutputPath, testerInstance.Output())
		})

		go testerWorker.Start()

		log.Info("yaja tester started ")

		// Wait for INT/TERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		for sig := range sigs {
			log.Infoln("received", sig)
			switch sig {
			case syscall.SIGTERM:
				log.Panic("terminated")
				os.Exit(0)
			case syscall.SIGQUIT:
				quitTester()
			case syscall.SIGHUP:
				quitTester()
			}
		}

	},
}

func quitTester() {
	testerWorker.Close()
	testerInstance.Close()
	srv.Close()

	file.SaveJSON(configTester.AccountsPath, db)
}

func init() {
	TesterCMD.Flags().StringVarP(&configPath, "config", "c", "yaja-tester.conf", "Path to configuration file")

}
