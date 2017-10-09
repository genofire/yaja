package cmd

import (
	"crypto/tls"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/model/config"

	"dev.sum7.eu/genofire/yaja/server"
	"github.com/genofire/golang-lib/worker"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var configPath string

var (
	configData      *config.Config
	state           *model.State
	statesaveWorker *worker.Worker
	srv             *server.Server
	certs           *tls.Config
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Runs the yaja server",
	Example: "yaja serve -c /etc/yaja.conf",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		configData, err = config.ReadConfigFile(configPath)
		if err != nil {
			log.Fatal("unable to load config file:", err)
		}

		state, err = model.ReadState(configData.StatePath)
		if err != nil {
			log.Warn("unable to load state file:", err)
		}

		statesaveWorker = worker.NewWorker(time.Minute, func() {
			model.SaveJSON(state, configData.StatePath)
			log.Info("save state to:", configData.StatePath)
		})

		m := autocert.Manager{
			Cache:  autocert.DirCache(configData.TLSDir),
			Prompt: autocert.AcceptTOS,
		}

		certs = &tls.Config{GetCertificate: m.GetCertificate}

		srv = &server.Server{
			TLSConfig:  certs,
			State:      state,
			PortClient: configData.PortClient,
			PortServer: configData.PortServer,
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

	model.SaveJSON(state, configData.StatePath)
}

func reload() {
	log.Info("start reloading...")
	configNewData, err := config.ReadConfigFile(configPath)
	if err != nil {
		log.Warn("unable to load config file:", err)
		return
	}

	if configNewData.StatePath != configData.StatePath {
		statesaveWorker.Close()
		statesaveWorker := worker.NewWorker(time.Minute, func() {
			model.SaveJSON(state, configNewData.StatePath)
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
		State:      state,
		PortClient: configNewData.PortClient,
		PortServer: configNewData.PortServer,
	}

	if configNewData.PortServer != configData.PortServer {
		restartServer = true
	}
	if configNewData.PortClient != configData.PortClient {
		restartServer = true
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
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "yaja.conf", "Path to configuration file")
}
