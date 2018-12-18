package main

import (
	"flag"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"moddns/app"
	"moddns/app/logger"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var VERSION = "0.0.1"

var (
	configFile string
	traceID    = uuid.New().String()
)

func init() {
	flag.StringVar(&configFile, "config", "", "configFile in (.json,.yaml,.toml)")
	flag.StringVar(&configFile, "c", "", "configFile in (.json,.yaml,.toml)")
}

func main() {
	flag.Parse()
	if configFile == "" {
		panic("Please use -c or -config local config")
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("Load local config error：" + err.Error())
	}

	var state int32 = 1
	ac := make(chan error)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGQUIT)

	httpHandler, closeHandle := app.Init(VERSION, traceID)

	go func() {
		logger.System(traceID).Infof("HTTP Server Starting , Port:[%s]", viper.GetString("http_port"))

		httpServer := &http.Server{
			Addr:           viper.GetString("http_addr"),
			Handler:        httpHandler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		ac <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-ac:
		if err != nil && atomic.LoadInt32(&state) == 1 {
		logger.System(traceID).Errorf("Listen HTTP server error:%s", err.Error())
	}
	case sig := <-sc:
		atomic.StoreInt32(&state, 0)
		logger.System(traceID).Infof("Get the exit signal[%s]", sig.String())

	}
	if closeHandle != nil {
		closeHandle()
	}

	os.Exit(int(atomic.LoadInt32(&state)))
}
