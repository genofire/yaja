package daemon

import (
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"dev.sum7.eu/genofire/golang-lib/file"
	"dev.sum7.eu/genofire/golang-lib/worker"
	"dev.sum7.eu/genofire/yaja/client"
	"dev.sum7.eu/genofire/yaja/daemon/tester"
	"dev.sum7.eu/genofire/yaja/messages"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var configTester = &tester.Config{}

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

		if err := file.ReadJSON(configTester.StatePath, db); err != nil {
			log.Warn("unable to load state file:", err)
		}

		statesaveWorker = worker.NewWorker(time.Minute, func() {
			file.SaveJSON(configTester.StatePath, db)
			log.Info("save state to:", configTester.StatePath)
		})

		// https server to handle acme (by letsencrypt)
		hs := &http.Server{
			Addr: configTester.Webserver,
		}
		if configTester.TLSDir != "" {
			m := autocert.Manager{
				Cache:  autocert.DirCache(configTester.TLSDir),
				Prompt: autocert.AcceptTOS,
			}
			hs.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
			go func(hs *http.Server) {
				if err := hs.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
					log.Errorf("webserver with addr %s: %s", hs.Addr, err)
				}
			}(hs)
		} else {
			go func(hs *http.Server) {
				if err := hs.ListenAndServe(); err != http.ErrServerClosed {
					log.Errorf("webserver with addr %s: %s", hs.Addr, err)
				}
			}(hs)
		}

		mainClient, err := client.NewClient(configTester.Client.JID, configTester.Client.Password)
		if err != nil {
			log.Fatal("unable to connect with main jabber client: ", err)
		}

		for _, admin := range configTester.Admins {
			mainClient.Out.Encode(&messages.MessageClient{
				From: mainClient.JID.Full(),
				To:   admin.Full(),
				Type: "chat",
				Body: "yaja tester starts",
			})
		}

		go statesaveWorker.Start()

		log.Infoln("yaja tester started ")

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
	srv.Close()
	statesaveWorker.Close()

	file.SaveJSON(configTester.StatePath, db)
}

func init() {
	TesterCMD.Flags().StringVarP(&configPath, "config", "c", "yaja-tester.conf", "Path to configuration file")

}
