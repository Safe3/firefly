package main

import (
	"fahi/pkg/config"
	"fahi/pkg/web"
	"fahi/pkg/wg"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	if os.Geteuid() != 0 {
		log.Fatal("please run firefly as root")
	}
	cfg, err := config.LoadOrCreate()
	if err != nil {
		log.Fatalf("failed to load or create config: %v", err)
	}

	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}
	log.SetLevel(logLevel)

	wgIface, err := wg.New(cfg)
	if err != nil {
		log.Fatalf("failed to init wireguard: %v", err)
	}
	defer wgIface.Close()

	err = wgIface.Create()
	if err != nil {
		log.Errorf("failed to create wireguard: %v", err)
		return
	}

	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	go func() {
		select {
		case <-termCh:
		}
		log.Info("shutdown signal received")
		wgIface.Close()
		os.Exit(1)
	}()

	err = web.Serve(cfg, wgIface)
	if err != nil {
		log.Errorf("failed to web server: %v", err)
	}
}
