package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joeshaw/envdecode"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type config struct {
	log  *logrus.Logger
	Addr string `env:"ADDR,default=localhost:8080"`
}

var (
	cfg *config
	log *logrus.Logger
)

func init() {
	cfg = &config{}
	envdecode.MustDecode(cfg)
	log = defaultLogger()
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting server...")
	serverConfig := httpServerConfig{
		log:    log,
		config: cfg,
	}
	server := startServer(ctx, serverConfig)

	log.Infof("System call: %+v. Shutting down...", <-sigChan)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Errorln("Failed to properly shutdown server")
	}

	cancel()
	<-shutdownCtx.Done()
	log.Infoln("Bye")
}

func defaultLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	log.SetLevel(logrus.TraceLevel)
	return log
}

type httpServerConfig struct {
	log    *logrus.Logger
	config *config
}

func startServer(ctx context.Context, hsc httpServerConfig) *http.Server {
	hh := healthHandler{
		log: hsc.log,
	}
	sh := sanitizeHandler{
		log: hsc.log,
	}
	r := mux.NewRouter()
	r.Handle("/health", &hh).Methods(http.MethodGet, http.MethodHead)
	r.Handle("/sanitize", &sh).Methods(http.MethodPost)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         hsc.config.Addr,
		Handler:      r,
	}

	go func() {
		hsc.log.Printf("Server listening on %s...", hsc.config.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hsc.log.WithError(err).Errorln("Failure with ListenAndServe")
		}
	}()

	return srv
}
