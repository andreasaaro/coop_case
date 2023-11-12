package app

import (
	"context"
	"fmt"
	"strconv"
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
	limitInt, err := strconv.Atoi(a.conf.MastodonConfig.Limit)
	if err != nil {
		return fmt.Errorf("error converting limit string to integer: %v", err)
	}

	blogCh := make(chan []*mastodon.MastodonData, limitInt)

	defer close(blogCh)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() error {
		defer wg.Done()
		if err := generateSourceStream(ctx, a.mastodon, blogCh); err != nil {
			return fmt.Errorf("error getting source data: %v", err)
		}
		return nil
	}()

	wg.Add(1)
	go func() error {
		defer wg.Done()
		if err := produce(ctx, blogCh, a.producer.Input(), a.conf.KafkaConfig); err != nil {
			return fmt.Errorf("error producing micro-blog messages to kafka: %v", err)
		}
		return nil
	}()

	wg.Add(1)
	go func() {
		defer logrus.Infof("exiting error handler")
		defer wg.Done()
		select {
		case <-ctx.Done():
			logrus.Infof("context canceled: %v", ctx.Err())
		case err := <-a.producer.Errors():
			logrus.Errorf("producer-error: %v", err)

		}
	}()

	wg.Wait()

	return nil
}
