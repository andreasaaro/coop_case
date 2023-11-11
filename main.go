package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"coop_case/app"

	"coop_case/config"
	"coop_case/kafka"
	"coop_case/mastodon"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Started app")
	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	go HandleSignals(ctx, cancel)

	mastodonClient, err := mastodon.NewClient(cfg.MastodonConfig)
	if err != nil {
		logrus.Fatalf("Failed to create mastodon client:%v", err)
	}

	kafkaProducer, err := kafka.NewSaramaProducer(cfg.KafkaConfig)
	if err != nil {
		logrus.Fatalf("Failed to create kafka producer:%v", err)
	}
	app := app.New(cfg, mastodonClient, kafkaProducer)

	if err := app.Run(ctx); err != nil && err != context.Canceled {
		logrus.Fatalf("App failed %v", err)
	}

	logrus.Info("Done")
}

func HandleSignals(ctx context.Context, cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	select {
	case <-ctx.Done():
	case <-signals:
		cancel()
	}
}
