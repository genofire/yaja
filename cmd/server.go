package cmd

import (
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/genofire/yaja/database"
	"github.com/genofire/yaja/model/config"

	"github.com/genofire/golang-lib/file"
	"github.com/genofire/golang-lib/worker"
	"github.com/genofire/yaja/server"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var configPath string

var (
	configData      = &config.Config{}
	db              = &database.State{}
	statesaveWorker *worker.Worker
	srv             *server.Server
	certs           *tls.Config
)

// serverCmd represents the serve command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Runs the yaja server",
	Example: "yaja serve -c /etc/yaja.conf",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		err = file.ReadTOML(configPath, configData)
		if err != nil {
			log.Fatal("unable to load config file:", err)
		}

		log.SetLevel(log.DebugLevel)

		err = file.ReadJSON(configData.StatePath, db)
		if err != nil {
			log.Warn("unable to load state file:", err)
		}

		statesaveWorker = worker.NewWorker(time.Minute, func() {
			file.SaveJSON(configData.StatePath, db)
			log.Info("save state to:", configData.StatePath)
		})

		m := autocert.Manager{
			Cache:  autocert.DirCache(configData.TLSDir),
			Prompt: autocert.AcceptTOS,
		}

		// https server to handle acme (by letsencrypt)
		httpServer := &http.Server{
			Addr:      ":https",
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		}
		go httpServer.ListenAndServeTLS("", "")

		srv = &server.Server{
			TLSManager: &m,
			Database:   db,
			ClientAddr: configData.Address.Client,
			ServerAddr: configData.Address.Server,
		}

		go statesaveWorker.Start()
		go srv.Start()

		log.Infoln("yaja started ")

		// Wait for INT/TERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1)
		for sig := range sigs {
			log.Infoln("received", sig)
			switch sig {
			case syscall.SIGTERM:
				log.Panic("terminated")
				os.Exit(0)
			case syscall.SIGQUIT:
				quit()
			case syscall.SIGHUP:
				quit()
			case syscall.SIGUSR1:
				reload()
			}
		}

	},
}

func quit() {
	srv.Close()
	statesaveWorker.Close()

	file.SaveJSON(configData.StatePath, db)
}

func reload() {
	log.Info("start reloading...")
	var configNewData *config.Config
	err := file.ReadTOML(configPath, configNewData)
	if err != nil {
		log.Warn("unable to load config file:", err)
		return
	}

	//TODO fetch changing address (to set restart)

	if configNewData.StatePath != configData.StatePath {
		statesaveWorker.Close()
		statesaveWorker := worker.NewWorker(time.Minute, func() {
			file.SaveJSON(configNewData.StatePath, db)
			log.Info("save state to:", configNewData.StatePath)
		})
		go statesaveWorker.Start()
	}

	restartServer := false

	if configNewData.TLSDir != configData.TLSDir {

		m := autocert.Manager{
			Cache:  autocert.DirCache(configData.TLSDir),
			Prompt: autocert.AcceptTOS,
		}

		certs = &tls.Config{GetCertificate: m.GetCertificate}
		restartServer = true
	}

	newServer := &server.Server{
		TLSConfig:  certs,
		Database:   db,
		ClientAddr: configNewData.Address.Client,
		ServerAddr: configNewData.Address.Server,
	}

	if restartServer {
		go srv.Start()
		//TODO should fetch new server error
		srv.Close()
		srv = newServer
	}

	configData = configNewData
	log.Info("reloaded")
}

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&configPath, "config", "c", "yaja.conf", "Path to configuration file")
}
