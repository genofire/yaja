package cmd

import (
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"dev.sum7.eu/genofire/yaja/database"
	"dev.sum7.eu/genofire/yaja/model/config"
	"dev.sum7.eu/genofire/yaja/server/extension"

	"dev.sum7.eu/genofire/golang-lib/file"
	"dev.sum7.eu/genofire/golang-lib/worker"
	"dev.sum7.eu/genofire/yaja/server"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var configPath string

var (
	configData       = &config.Config{}
	db               = &database.State{}
	statesaveWorker  *worker.Worker
	srv              *server.Server
	certs            *tls.Config
	extensionsClient extension.Extensions
	extensionsServer extension.Extensions
)

// serverCmd represents the serve command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Runs the yaja server",
	Example: "yaja serve -c /etc/yaja.conf",
	Run: func(cmd *cobra.Command, args []string) {

		if err := file.ReadTOML(configPath, configData); err != nil {
			log.Fatal("unable to load config file:", err)
		}

		log.SetLevel(configData.Logging.Level)

		if err := file.ReadJSON(configData.StatePath, db); err != nil {
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
		for _, addr := range configData.Address.Webserver {
			hs := &http.Server{
				Addr:      addr,
				TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
			}
			go func(hs *http.Server, addr string) {
				if err := hs.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
					log.Errorf("webserver with addr %s: %s", addr, err)
				}
			}(hs, addr)
		}

		srv = &server.Server{
			TLSManager:       &m,
			Database:         db,
			ClientAddr:       configData.Address.Client,
			ServerAddr:       configData.Address.Server,
			LoggingClient:    configData.Logging.LevelClient,
			LoggingServer:    configData.Logging.LevelServer,
			RegisterEnable:   configData.Register.Enable,
			RegisterDomains:  configData.Register.Domains,
			ExtensionsServer: extensionsServer,
			ExtensionsClient: extensionsClient,
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

	if err := file.ReadTOML(configPath, configNewData); err != nil {
		log.Warn("unable to load config file:", err)
		return
	}
	log.SetLevel(configNewData.Logging.Level)
	srv.LoggingClient = configNewData.Logging.LevelClient
	srv.LoggingServer = configNewData.Logging.LevelServer
	srv.RegisterEnable = configNewData.Register.Enable
	srv.RegisterDomains = configNewData.Register.Domains

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
	if restartServer {
		newServer := &server.Server{
			TLSConfig:        certs,
			Database:         db,
			ClientAddr:       configNewData.Address.Client,
			ServerAddr:       configNewData.Address.Server,
			LoggingClient:    configNewData.Logging.LevelClient,
			RegisterEnable:   configNewData.Register.Enable,
			RegisterDomains:  configNewData.Register.Domains,
			ExtensionsServer: extensionsServer,
			ExtensionsClient: extensionsClient,
		}
		log.Warn("reloading need a restart:")
		go newServer.Start()
		//TODO should fetch new server error
		srv.Close()
		srv = newServer
	}

	configData = configNewData
	log.Info("reloaded")
}

func init() {
	extensionsClient = append(extensionsClient,
		&extension.Message{},
		&extension.Presence{},
		extension.IQExtensions{
			&extension.IQPrivate{},
			&extension.IQPing{},
			&extension.IQLast{},
			&extension.IQDisco{Database: db},
			&extension.IQRoster{Database: db},
			&extension.IQExtensionDiscovery{GetSpaces: func() []string {
				return extensionsClient.Spaces()
			}},
		})

	extensionsServer = append(extensionsServer,
		extension.IQExtensions{
			&extension.IQPing{},
		})

	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&configPath, "config", "c", "yaja.conf", "Path to configuration file")

}
