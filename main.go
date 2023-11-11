package main

import (
	"context"
	"fmt"

	"coop_case/app"

	"coop_case/config"
	"coop_case/kafka"
	"coop_case/mastodon"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("started app")
	cfg := config.NewConfig()
	// logrus.Infoln(cfg.PrintConfig(AppName, Version))

	ctx, cancel := context.WithCancel(context.Background())

	mastodonClient, err := mastodon.NewClient(cfg.MastodonConfig)
	if err != nil {
		logrus.Info("unable to create mastodon client")
	}

	kafkaProducer, err := kafka.NewSaramaProducer(cfg.KafkaConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	app := app.New(cfg, mastodonClient, kafkaProducer)

	if err := app.Run(ctx); err != nil && err != context.Canceled {
		logrus.Fatalf("App failed %v", err)
	}

	logrus.Info("Done")

	fmt.Println(cancel)
}
