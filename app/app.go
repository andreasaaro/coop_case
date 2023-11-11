package app

import (
	"context"
	"fmt"
	"sync"

	"coop_case/config"
	"coop_case/kafka"
	"coop_case/mastodon"

	"github.com/sirupsen/logrus"
)

type App interface {
	Run(context.Context) error
}

type app struct {
	conf     *config.Config
	mastodon mastodon.Mastodon
	producer kafka.Producer
}

func New(conf *config.Config, mastodon mastodon.Mastodon, producer kafka.Producer) App {
	return &app{
		conf:     conf,
		mastodon: mastodon,
		producer: producer,
	}
}

func (a *app) Run(ctx context.Context) error {
	blogCh := make(chan []*mastodon.MastodonData)

	defer close(blogCh)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := generateSourceStream(ctx, a.mastodon, blogCh); err != nil {
			fmt.Println("error gen source")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := produce(ctx, blogCh, a.producer.Input(), a.conf.KafkaConfig); err != nil {
			fmt.Println("error produce source")
		}
	}()

	wg.Add(1)
	go func() {
		defer logrus.Infof("exiting error handler")
		defer wg.Done()
		select {
		case <-ctx.Done():
			logrus.Infof("%v", ctx.Err())
		case err := <-a.producer.Errors():
			logrus.Infof("Producer-error: %v", err)

		}
	}()

	wg.Wait()

	return nil
}
