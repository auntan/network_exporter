package main

import (
	"context"
	"flag"
	"github.com/auntan/network_exporter/internal/config"
	"github.com/auntan/network_exporter/internal/graceful"
	"github.com/auntan/network_exporter/internal/http_server"
	"github.com/auntan/network_exporter/internal/log"
	"github.com/auntan/network_exporter/internal/metrics"
	"github.com/auntan/network_exporter/internal/pinger"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

var configPath = flag.String("config", "configs/config.yaml", "configuration file")

func main() {
	flag.Parse()

	log.InitLoggerDefault()

	conf, err := config.Load(*configPath)
	if err != nil {
		zap.S().Fatalf("load config error: %v", err)
	}

	err = log.InitLogger(conf)
	if err != nil {
		return
	}
	defer zap.S().Sync()
	zap.S().Info("started")

	metrics.Initialize(&metrics.Config{
		RTTHistogramBuckets: conf.HistogramBuckets,
	})

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = graceful.Run(ctx, map[string]func(context.Context) error{
		"http": http_server.New(&http_server.Config{
			Port: conf.HttpPort,
		}).Run,
		"pings": pinger.New(conf).Run,
	})
	if err != nil {
		zap.S().Errorf("run error: %v", err)
	}

	zap.S().Info("bye")
}
